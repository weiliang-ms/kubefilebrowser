package terminal

import (
	"encoding/json"
	"fmt"
	"kubefilebrowser/configs"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/logs"
	"os"
	"path/filepath"
	"time"
)

const header = `{version":2,"width":153,"height":44,"timestamp":%d,"env":{"SHELL":"%s","TERM":"xterm-256color"}}`

type ReplyRecorder struct {
	sessionID     string
	shell         string
	timeStartNano int64
	absFilePath   string
	AbsGzFilePath string
	file          *os.File
}

func NewReplyRecord(sid, shell string) (recorder ReplyRecorder) {
	recorder = ReplyRecorder{sessionID: sid, shell: shell}
	recorder.initial()
	return recorder
}

func (r *ReplyRecorder) initial() {
	r.prepare()
}

func (r *ReplyRecorder) Record(b []byte) {
	if len(b) > 0 {
		delta := float64(time.Now().UnixNano()-r.timeStartNano) / 1000 / 1000 / 1000
		data, _ := json.Marshal(string(b))
		_, err := r.file.WriteString(fmt.Sprintf(`["%f","o",%s]`, delta, data))
		if err != nil {
			logs.Error(fmt.Sprintf("Session %s write replay to file failed: %s", r.sessionID, err))
		}
		_, _ = r.file.Write([]byte("\r\n"))
	}
}

func (r *ReplyRecorder) prepare() {
	sessionID := r.sessionID
	rootPath := configs.Config.RootPath
	today := time.Now().UTC().Format("2006-01-02")
	gzFileName := sessionID + ".replay.gz"
	replayDir := filepath.Join(rootPath, "data", "replays", today)

	r.absFilePath = filepath.Join(replayDir, sessionID)
	r.AbsGzFilePath = filepath.Join(replayDir, gzFileName)
	r.timeStartNano = time.Now().UnixNano()

	err := utils.EnsureDirExist(replayDir)
	if err != nil {
		logs.Error(fmt.Sprintf("Create dir %s error: %s\n", replayDir, err))
		return
	}

	logs.Info(fmt.Sprintf("Session: %s, Replay file path: %s", r.sessionID, r.absFilePath))
	r.file, err = os.Create(r.absFilePath)
	if err != nil {
		logs.Error(fmt.Sprintf("Create file %s error: %s", r.absFilePath, err))
	}
	_, _ = r.file.WriteString(fmt.Sprintf(header, time.Now().Nanosecond(), r.shell))
	_, _ = r.file.Write([]byte("\r\n"))
}

func (r *ReplyRecorder) End() {
	delta := float64(time.Now().UnixNano()-r.timeStartNano) / 1000 / 1000 / 1000
	_, _ = r.file.WriteString(fmt.Sprintf(`["%f":"","%f":""]`, delta, 0.0))
	_ = r.file.Close()
	if !utils.FileOrPathExist(r.AbsGzFilePath) {
		logs.Info("Compress replay file: ", r.absFilePath)
		_ = utils.GzipCompressFile(r.absFilePath, r.AbsGzFilePath)
		_ = os.Remove(r.absFilePath)
	}
	//go r.UploadGzipFile()
}
