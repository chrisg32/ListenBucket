namespace ListenBucket.Data.Models;

public class Podcast
{
    public int Id { get; set; }
    public string Title { get; set; }
    public string Description { get; set; }
    public string ThumbnailUrl { get; set; }
    // public string Author { get; set; }
    public string? Source { get; set; }
    public string? YouTubeId { get; set; }
    public string? Byline { get; set; }
    public string? BylineLogoUrl { get; set; }
    public DateTimeOffset? LastEpisodeDate { get; set; }
    public int? EpisodeCount { get; set; }
}