using YoutubeDLSharp;
using YoutubeDLSharp.Options;

namespace ListenBucket.Services;

public static class DownloadService
{
    public static async Task<string?> DownloadAsync(string url, string directory, string fileName)
    {
        var downloader = new YoutubeDL
        {
            OutputFolder = directory,
            YoutubeDLPath = "yt-dlp",
            FFmpegPath = "ffmpeg",
        };

        var res = await downloader.RunAudioDownload(
            url,
            AudioConversionFormat.Mp3,
            overrideOptions: new OptionSet
            {
                Output = $"{fileName}.%(ext)s"
            }
        );
        
        res.EnsureSuccess();

        return res.Data;
    }
}