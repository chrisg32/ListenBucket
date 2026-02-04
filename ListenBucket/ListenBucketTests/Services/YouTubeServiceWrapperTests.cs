using ListenBucket.Services;

namespace ListenBucketTests.Services;

public class YouTubeServiceWrapperTests
{
    [Fact]
    public async Task Test_GetMetaData_Channel_UserHandle()
    {
        // Arrange
        var youtubeService = new YouTubeServiceWrapper();
        var channelUrl = "https://www.youtube.com/@TalkShiftPod";

        // Act
        var metaData = await youtubeService.GetMetaData(channelUrl, true);

        // Assert
        Assert.NotNull(metaData);
        Assert.Equal("Talk Shift", metaData.Title);
        Assert.Equal(YouTubeServiceWrapper.YouTubeSourceType.Channel, metaData.SourceType);
        Assert.NotNull(metaData.Children);
        Assert.NotEmpty(metaData.Children);
    }
    
    [Fact]
    public async Task Test_GetMetaData_Video()
    {
        // Arrange
        var youtubeService = new YouTubeServiceWrapper();
        var videoUrl = "https://www.youtube.com/watch?v=dQw4w9WgXcQ";

        // Act
        var metaData = await youtubeService.GetMetaData(videoUrl, true);

        // Assert
        Assert.NotNull(metaData);
        Assert.Equal("Rick Astley - Never Gonna Give You Up (Official Music Video)", metaData.Title); // Replace with actual title
        Assert.Equal(YouTubeServiceWrapper.YouTubeSourceType.Video, metaData.SourceType);
    }

    [Fact]
    public async Task Test_GetMetaData_Playlist()
    {
        // Arrange
        var youtubeService = new YouTubeServiceWrapper();
        var playlistUrl = "https://www.youtube.com/playlist?list=PL8mG-RkN2uTwpDRywu9O2JaRBeOYgNGNa";

        // Act
        var metaData = await youtubeService.GetMetaData(playlistUrl, true);

        // Assert
        Assert.NotNull(metaData);
        Assert.Equal("Live Stream", metaData.Title);
        Assert.Equal("Enjoy the mulling and crazy almost unscripted events at Linus Tech Tips. Usually happens on Friday", metaData.Description);
        Assert.Equal(YouTubeServiceWrapper.YouTubeSourceType.Playlist, metaData.SourceType);
        Assert.NotNull(metaData.Children);
        Assert.NotEmpty(metaData.Children);
    }

    [Fact]
    public async Task Test_GetMetaData_Channel()
    {
        // Arrange
        var youtubeService = new YouTubeServiceWrapper();
        var channelUrl = "https://www.youtube.com/channel/UC_x5XG1OV2P6uZZ5FSM9Ttw";

        // Act
        var metaData = await youtubeService.GetMetaData(channelUrl, true);

        // Assert
        Assert.NotNull(metaData);
        Assert.Equal("Google for Developers", metaData.Title);
        Assert.Equal(YouTubeServiceWrapper.YouTubeSourceType.Channel, metaData.SourceType);
    }
}