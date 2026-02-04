using Google.Apis.Services;
using Google.Apis.YouTube.v3;

namespace ListenBucket.Services;

public class YouTubeServiceWrapper
{
    private const string ApiKey = "AIzaSyDRH8I2fEI4pWI-WIif-iamw2J_bHTB5Rs";
    
    private readonly YouTubeService _youtubeService = new(new BaseClientService.Initializer()
    {
        ApiKey = ApiKey
    });
    
    public async Task<YouTubeMetadata> GetMetaData(string url, bool includeChildren)
    {
        try
        {
            var uri = new Uri(url);
            var query = System.Web.HttpUtility.ParseQueryString(uri.Query);

            if (query["v"] != null) // Video
            {
                var videoId = query["v"];
                return await GetVideoMetadata(videoId);
            }

            if (query["list"] != null) // Playlist
            {
                var playlistId = query["list"];
                return await GetPlaylistMetadata(playlistId, includeChildren);
            }

            if (uri.AbsolutePath.StartsWith("/channel/")) // Channel
            {
                var channelId = uri.AbsolutePath.Split('/').Last();
                return await GetChannelMetadata(channelId, includeChildren);
            }
            
            if (uri.AbsolutePath.StartsWith("/@")) // Channel with handle
            {
                var channelHandle = uri.AbsolutePath.TrimStart('/');
                return await GetChannelMetadataByHandle(channelHandle, includeChildren);
            }

            throw new ArgumentException("Invalid YouTube URL");
        }
        catch (Exception e)
        {
            Console.WriteLine(e);
            throw;
        }
    }

    private async Task<YouTubeMetadata> GetVideoMetadata(string? videoId)
    {
        var request = _youtubeService.Videos.List("snippet");
        request.Id = videoId;

        var response = await request.ExecuteAsync();
        var video = response.Items?.FirstOrDefault();

        if (video == null)
        {
            throw new Exception("Video not found");
        }

        return new YouTubeMetadata
        {
            Title = video.Snippet.Title,
            Description = video.Snippet.Description,
            PublishedAt = video.Snippet.PublishedAtDateTimeOffset ?? DateTimeOffset.Now,
            ThumbnailUrl = video.Snippet.Thumbnails.Standard?.Url ?? video.Snippet.Thumbnails.Default__.Url,
            YouTubeId = video.Id,
            SourceType = YouTubeSourceType.Video
        };
    }

    private async Task<YouTubeMetadata> GetPlaylistMetadata(string? playlistId, bool includeChildren = true)
    {
        var request = _youtubeService.Playlists.List("snippet");
        request.Id = playlistId;

        var response = await request.ExecuteAsync();
        var playlist = response.Items?.FirstOrDefault();

        if (playlist == null)
        {
            throw new Exception("Playlist not found");
        }

        var metadata = new YouTubeMetadata
        {
            Title = playlist.Snippet.Title,
            Description = playlist.Snippet.Description,
            PublishedAt = playlist.Snippet.PublishedAtDateTimeOffset ?? DateTimeOffset.Now,
            ThumbnailUrl = playlist.Snippet.Thumbnails.Standard?.Url ?? playlist.Snippet.Thumbnails.Default__.Url,
            YouTubeId = playlist.Id,
            SourceType = YouTubeSourceType.Playlist
        };

        if (includeChildren)
        {
            metadata.Children = await GetAllVideosInPlaylist(playlistId);
        }

        return metadata;
    }

    private async Task<List<YouTubeMetadata>> GetAllVideosInPlaylist(string? playlistId)
    {
        var videos = new List<YouTubeMetadata>();
        if (playlistId == null) return videos;
        string? nextPageToken = null;

        do
        {
            var playlistItemsRequest = _youtubeService.PlaylistItems.List("snippet");
            playlistItemsRequest.PlaylistId = playlistId;
            playlistItemsRequest.MaxResults = 50;
            playlistItemsRequest.PageToken = nextPageToken;

            var playlistItemsResponse = await playlistItemsRequest.ExecuteAsync();
            nextPageToken = playlistItemsResponse.NextPageToken;

            videos.AddRange(playlistItemsResponse.Items.Select(item => new YouTubeMetadata
            {
                Title = item.Snippet.Title,
                Description = item.Snippet.Description,
                PublishedAt = item.Snippet.PublishedAtDateTimeOffset ?? DateTimeOffset.Now,
                ThumbnailUrl = item.Snippet.Thumbnails.Standard?.Url ?? item.Snippet.Thumbnails.Default__.Url,
                YouTubeId = item.Snippet.ResourceId.VideoId,
                SourceType = YouTubeSourceType.Video
            }));
        } while (!string.IsNullOrEmpty(nextPageToken));

        return videos;
    }
    
    private async Task<YouTubeMetadata> GetChannelMetadataByHandle(string channelHandle, bool includeChildren)
    {
        // Use the search endpoint to find the channel by handle
        var searchRequest = _youtubeService.Search.List("snippet");
        searchRequest.Q = channelHandle;
        searchRequest.Type = "channel";
        searchRequest.MaxResults = 1;

        var searchResponse = await searchRequest.ExecuteAsync();
        var channel = searchResponse.Items?.FirstOrDefault();

        if (channel == null)
        {
            Console.WriteLine(channelHandle);
            throw new Exception("Channel not found");
        }

        // Use the channel ID to fetch detailed metadata
        var channelId = channel.Id.ChannelId;
        return await GetChannelMetadata(channelId, includeChildren);
    }

    private async Task<YouTubeMetadata> GetChannelMetadata(string channelId, bool includeChildren)
    {
        var request = _youtubeService.Channels.List("snippet");
        request.Id = channelId;

        var response = await request.ExecuteAsync();
        var channel = response.Items?.FirstOrDefault();

        if (channel == null)
        {
            throw new Exception("Channel not found");
        }

        var metadata = new YouTubeMetadata
        {
            Title = channel.Snippet.Title,
            Description = channel.Snippet.Description,
            PublishedAt = channel.Snippet.PublishedAtDateTimeOffset ?? DateTimeOffset.Now,
            ThumbnailUrl = channel.Snippet.Thumbnails.Standard?.Url ?? channel.Snippet.Thumbnails.Default__.Url,
            YouTubeId = channel.Id,
            SourceType = YouTubeSourceType.Channel
        };

        if (includeChildren)
        {
            metadata.Children = await GetAllVideosByChannel(channelId);
        }

        return metadata;
    }

    private async Task<List<YouTubeMetadata>> GetAllVideosByChannel(string channelId)
    {
        var videos = new List<YouTubeMetadata>();
        string? nextPageToken = null;

        do
        {
            var searchRequest = _youtubeService.Search.List("snippet");
            searchRequest.ChannelId = channelId;
            searchRequest.Type = "video";
            searchRequest.MaxResults = 50;
            searchRequest.PageToken = nextPageToken;

            var searchResponse = await searchRequest.ExecuteAsync();
            nextPageToken = searchResponse.NextPageToken;

            videos.AddRange(searchResponse.Items.Select(video => new YouTubeMetadata
            {
                Title = video.Snippet.Title,
                Description = video.Snippet.Description,
                PublishedAt = video.Snippet.PublishedAtDateTimeOffset ?? DateTimeOffset.Now,
                ThumbnailUrl = video.Snippet.Thumbnails.Standard?.Url ?? video.Snippet.Thumbnails.Default__.Url,
                YouTubeId = video.Id.VideoId,
                SourceType = YouTubeSourceType.Video
            }));
        } while (!string.IsNullOrEmpty(nextPageToken));

        return videos;
    }

    public class YouTubeMetadata
    {
        public string Title { get; set; }
        public string Description { get; set; }
        public string ThumbnailUrl { get; set; }
        public DateTimeOffset PublishedAt { get; set; }
        public string YouTubeId { get; set; }
        public YouTubeSourceType SourceType { get; set; }
        public List<YouTubeMetadata>? Children { get; set; }
    }
    
    public enum YouTubeSourceType
    {
        Video,
        Playlist,
        Channel
    }
}