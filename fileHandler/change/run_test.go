package fileHandler

import (
	"github.com/creationtime/lib-go/fileHandler"
	"testing"
)

func Test_voice(t *testing.T) {
	wavFile := "/Users/tim/Desktop/images/voice/1614324496966669544blob" // 需要转换的wav文件
	mp3File := "/Users/tim/Desktop/images/voice/1.mp3"                   // 转换后mp3文件存放路径
	if err := VoiceHandle(wavFile, mp3File); err != nil {
		t.Logf("===> old err:%v", err.Error())
	}
}

func Test_voice2(t *testing.T) {
	wavFile := "/Users/tim/Music/网易云音乐/静心/般禅梵唱妙音组 - 大准提咒.flac" // 需要转换的文件
	mp3File := "/Users/tim/Music/网易云音乐/静心/般禅梵唱妙音组 - 大准提咒.mp3"  // 转换后mp3文件存放路径
	if err := VoiceChange(wavFile, mp3File); err != nil {
		t.Logf("===> old err:%v", err.Error())
	}
}

func Test_voice3(t *testing.T) {
	if list, err := fileHandler.GetFile("/Users/tim/Music/网易云音乐/静心"); err != nil {
		t.Logf("===> old err: %v", err.Error())
	} else {
		if len(list) > 0 {
			for _, v := range list {
				t.Logf("====>文件: %v", v)
				//读取文件
			}
		}
		//wavFile := "/Users/tim/Music/网易云音乐/静心/般禅梵唱妙音组 - 大准提咒.flac" // 需要转换的文件
		//mp3File := "/Users/tim/Music/网易云音乐/静心/般禅梵唱妙音组 - 大准提咒.mp3"  // 转换后mp3文件存放路径
		//if err := VoiceChange(wavFile, mp3File); err != nil {
		//	t.Logf("===> old err:%v", err.Error())
		//}
	}
}
