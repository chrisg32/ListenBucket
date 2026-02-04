namespace ListenBucket.Data.Models;

public class Item
{
    public int Id { get; set; }
    public int PodcastId { get; set; }
    public string? Source { get; set; }
    public string YouTubeId { get; set; }
    public string Title { get; set; }
    public string Description { get; set; }
    public DateTimeOffset PublishedAt { get; set; }
    public string ThumbnailUrl { get; set; }
    public bool Downloaded { get; set; }
    public TimeSpan? Duration { get; set; }
    public DateTimeOffset? AddedAt { get; set; }
}