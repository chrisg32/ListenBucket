using System.Net;
using System.ServiceModel.Syndication;
using System.Xml;
using ListenBucket.Data;
using ListenBucket.Data.Models;
using ListenBucket.Services.Podcast;
using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;

namespace ListenBucket.Controllers;

public class FeedController : Controller
{
    private const int ListenLaterPodcastId = 1;
    private readonly ApplicationDbContext _context;

    public FeedController(ApplicationDbContext context)
    {
        _context = context;
    }



    [HttpGet("/feed/listen-later")]
    [HttpGet("/feed/listenLater")]
    public async Task GetListenLater(CancellationToken cancellationToken)
    {
        await GetFeed(cancellationToken, ListenLaterPodcastId);
    }
    
    [HttpGet("/feed/{podcastId}")]
    public async Task GetFeed(CancellationToken cancellationToken, int podcastId)
    {
        var hostName = Environment.GetEnvironmentVariable("HOSTNAME") ?? "localhost";
        //get podcast from database
        var podcast = await _context.Podcasts.FirstOrDefaultAsync(p => p.Id == podcastId, cancellationToken);
        if(podcast == null) 
        {
            if (podcastId == ListenLaterPodcastId)
            {
                podcast = new Podcast
                {
                    Id = ListenLaterPodcastId,
                    Title = "Listen Later",
                    Description = "A collection of items to listen to later",
                    ThumbnailUrl = "/i/listen-later.png",
                };
            }
            else
            {
                Response.StatusCode = (int)HttpStatusCode.NotFound;
                return;
            }
        }
        
        //create feed
        var feed = new PodcastSyndicationFeed
        {
            Title = podcast.Title,
            Description = podcast.Description,
            // Author = string.IsNullOrWhiteSpace(podcast.Author) ? "ListenBucket" : podcast.Author,
            Author = "ListenBucket",
            Categories = [PodcastCategory.Technology],
            AlternativeLink = new Uri("http://" +  hostName + "/podcast/" + podcast.Id),
            LastUpdatedTime = DateTimeOffset.Now,
            Language = "en-US",
            Explicit = false,
            PodcastType = PodcastType.Episodic,
            ImageUrl = new Uri("http://" +  hostName + "/i/listen-later.png")
        };
        
        // add podcast items
        var items = await _context.Items.Where(i => i.PodcastId == podcastId).ToListAsync(cancellationToken);
        var feedItems = items.OrderByDescending(i => i.AddedAt)
            .Select(item => new PodcastItem
            {
                Title = item.Title,
                Description = item.Description,
                PublishedDate = item.PublishedAt,
                Id = item.YouTubeId,
                LastUpdatedTime = item.PublishedAt,
                MediaUrl = new Uri($"http://{hostName}/m/{item.YouTubeId}.mp3"),
                ImageUrl = string.IsNullOrWhiteSpace(item.ThumbnailUrl) ? null : new Uri(item.ThumbnailUrl)
            })
            .ToList();
        
        feed.Items = feedItems;
        
        Response.StatusCode = (int)HttpStatusCode.OK;
        Response.ContentType = "application/rss+xml; charset=utf-8";

        await feed.WriteToStreamAsync(Response.Body);
    }
}