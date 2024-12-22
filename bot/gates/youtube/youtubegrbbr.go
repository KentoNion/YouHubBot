package youtube

import (
	"io"
	"log"
	"os"

	youtube "github.com/kkdai/youtube/v2"
)

// YoutubeService структура для работы с YouTube API
type YoutubeService struct {
	client *youtube.Client
}

// NewYoutubeService создает новую службу для работы с YouTube
func NewYoutubeService() *YoutubeService {
	return &YoutubeService{
		client: &youtube.Client{},
	}
}

// GetVideoInfo получает информацию о видео по ссылке
func (s *YoutubeService) GetVideoInfo(videoURL string) (*youtube.Video, error) {
	video, err := s.client.GetVideo(videoURL)
	if err != nil {
		log.Printf("Ошибка при получении видео: %v", err)
		return nil, err
	}
	return video, nil
}

// DownloadVideo загружает видео по ссылке
func (s *YoutubeService) DownloadVideo(videoURL string, outputPath string) error {
	video, err := s.GetVideoInfo(videoURL)
	if err != nil {
		return err
	}

	stream, _, err := s.client.GetStream(video, &video.Formats[0])
	if err != nil {
		log.Printf("Ошибка при загрузке видео: %v", err)
		return err
	}

	// Создание выходного файла
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Запись потока в файл
	_, err = io.Copy(file, stream)
	if err != nil {
		log.Printf("Ошибка записи видео: %v", err)
		return err
	}

	return nil
}
