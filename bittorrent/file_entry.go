package bittorrent

import (
	"io"
	"os"
	"path/filepath"

	"github.com/anacrolix/missinggo"
	gotorrent "github.com/anacrolix/torrent"

	"github.com/elgatito/elementum/broadcast"
	"github.com/elgatito/elementum/util"
	"github.com/elgatito/elementum/xbmc"
)

type SeekableContent interface {
	io.ReadSeeker
	io.Closer
}

type FileEntry struct {
	*Torrent
	*gotorrent.File
	*Reader

	rs                 io.ReadSeeker
	libraryBroadcaster *broadcast.Broadcaster
}

func (e *FileEntry) Seek(offset int64, whence int) (int64, error) {
	return e.Reader.Seek(offset+e.File.Offset(), whence)
}

func NewFileReader(t *Torrent, f *gotorrent.File, sequential bool) (*FileEntry, error) {
	var reader *Reader
	if t.IsPlaying {
		reader = t.GetActiveReader(f)
	} else {
		reader = t.NewReader(f)

		if sequential {
			reader.SetResponsive()
		}

		tfsLog.Infof("Setting readahead for reader %d as %d", reader.id, t.Service.GetReadaheadSize())
		reader.SetReadahead(t.Service.GetReadaheadSize())
	}

	if _, err := reader.Seek(f.Offset(), os.SEEK_SET); err != nil {
		return nil, err
	}

	rs := missinggo.NewSectionReadSeeker(reader, f.Offset(), f.Length())

	entry := &FileEntry{
		Torrent: t,
		File:    f,
		Reader:  reader,
		rs:      rs,

		libraryBroadcaster: broadcast.LocalBroadcasters[broadcast.WATCHED],
	}

	entry.setSubtitles()

	return entry, nil
}

func (e *FileEntry) setSubtitles() {
	filePath := e.File.Path()
	extension := filepath.Ext(filePath)

	if extension != ".srt" {
		srtPath := filePath[0:len(filePath)-len(extension)] + ".srt"
		files := e.Torrent.Files()

		for _, f := range files {
			if f.Path() == srtPath {
				xbmc.PlayerSetSubtitles(util.GetHTTPHost() + "/files/" + srtPath)
				return
			}
		}
	}
}

func (e *FileEntry) Close() error {
	tfsLog.Info("Closing file...")

	if item := e.Torrent.GetPlayingItem(); item != nil {
		e.libraryBroadcaster.Broadcast(item)
	}

	return e.Reader.Close()
}
