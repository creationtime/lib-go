package fileHandler

import (
	"os/exec"

	"github.com/xfrr/goffmpeg/transcoder"
)

//需先安装linux的ffmpeg命令
//mac下安装 brew install ffmpeg
func VoiceHandle(inFile, outFile string) error {
	//wav_file := "/root/input.wav"  // 需要转换的wav文件
	//mp3_file := "/root/output.mp3" // 转换后mp3文件存放路径
	//cmd := exec.Command("ffmpeg", inFile, outFile)
	cmd := exec.Command("ffmpeg", "-ss", "0", "-t", "20", "-i", inFile, "-f", "s16le", "-acodec", "pcm_s16le", "-b:a", "16", "-ar", "8000", "-ac", "1", outFile)
	//cmd := exec.Command("lame", inFile, outFile)
	err := cmd.Run()
	if err != nil {
		return err
	}
	// wav转mp3成功后，如有必要则可删除wav原文件
	//os.Remove(wav_file)
	return nil
}

func VoiceChange(inFile, outFile string) error {
	trans := new(transcoder.Transcoder)

	err := trans.Initialize(inFile, outFile)
	if err != nil {
		return err
	}
	done := trans.Run(false)
	err = <-done
	if err != nil {
		return err
	}
	// wav转mp3成功后，如有必要则可删除wav原文件
	//_ = os.Remove(inFile)
	return nil
}
