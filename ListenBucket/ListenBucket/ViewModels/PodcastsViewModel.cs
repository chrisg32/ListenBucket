using System.Collections.ObjectModel;
using ListenBucket.Data;
using ListenBucket.Data.Models;
using ListenBucket.Services;
using Microsoft.EntityFrameworkCore;

namespace ListenBucket.ViewModels;

public class PodcastsViewModel(ApplicationDbContext context)
{
    public ObservableCollection<Podcast> Podcasts { get; } = new();
    public string? Url { get; set; } = "https://www.youtube.com/@TalkShiftPod";

    public async Task Add()
    {
        var youtubeService = new YouTubeServiceWrapper();
        var meta = await youtubeService.GetMetaData(Url, true);
        if (meta.SourceType == YouTubeServiceWrapper.YouTubeSourceType.Video) return;
        
        var podcast = new Podcast
        {
            Title = meta.Title,
            Description = meta.Description,
            ThumbnailUrl = meta.ThumbnailUrl,
            YouTubeId = meta.YouTubeId,
            LastEpisodeDate = meta.PublishedAt,
            EpisodeCount = meta.Children?.Count ?? 0
        };
        context.Podcasts.Add(podcast);

        if (meta.Children?.Any() == true)
        {
            var items = meta.Children.Select(c => new Item
            {
                Source = c.YouTubeId,
                YouTubeId = c.YouTubeId,
                Title = c.Title,
                Description = c.Description,
                PublishedAt = c.PublishedAt,
                ThumbnailUrl = c.ThumbnailUrl,
                PodcastId = podcast.Id
            }).ToList();
            context.Items.AddRange(items);
        }
        
        await context.SaveChangesAsync();
        Podcasts.Add(podcast);
        Url = null;
    }
    
    public async Task Load()
    {
        var items = await context.Podcasts.Where(p => p.Id != 1).ToListAsync();
        foreach (var item in items.OrderByDescending(i => i.Title))
        {
            Podcasts.Add(item);
        }
    }
}