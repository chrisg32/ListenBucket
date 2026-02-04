using ListenBucket.Data;
using ListenBucket.Data.Models;
using Microsoft.EntityFrameworkCore;
using System.Collections.ObjectModel;
using System.Threading.Tasks;
using ListenBucket.Services;

namespace ListenBucket.ViewModels;

public class ListenLaterViewModel(ApplicationDbContext context, JobService jobService)
{
    public ObservableCollection<Item> Items { get; } = new();

    public string Url { get; set; }

    public string FeedUrl
    {
        get
        {
            var hostName = Environment.GetEnvironmentVariable("HOSTNAME") ?? "localhost";
            if (Uri.TryCreate($"http://{hostName}", UriKind.Absolute, out var uri))
            {
                var builder = new UriBuilder(uri)
                {
                    Path = "/feed/listen-later"
                };
                return builder.ToString();
            }
            return string.Empty;
        }
    }

    public async Task Load()
    {
        var items = await context.Items.Where(i => i.PodcastId == 1).ToListAsync();
        foreach (var item in items.OrderByDescending(i => i.AddedAt))
        {
            Items.Add(item);
        }
    }

    public async Task Add()
    {
        var youtubeService = new YouTubeServiceWrapper();
        var meta = await youtubeService.GetMetaData(Url, true);
        //TODO handle meta data type is not a video, if it is a playlist or channel, we should add all items in the playlist or channel
        var item = new Item
        {
            Source = Url,
            YouTubeId = meta.YouTubeId,
            Title = meta.Title,
            Description = meta.Description,
            PublishedAt = meta.PublishedAt,
            ThumbnailUrl = meta.ThumbnailUrl,
            AddedAt = DateTimeOffset.Now,
            PodcastId = 1
        };
        context.Items.Add(item);
        await context.SaveChangesAsync();
        
        jobService.QueueDownload(item.Id);
        
        Items.Add(item);
    }
}