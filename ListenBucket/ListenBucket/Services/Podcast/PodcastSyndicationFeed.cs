using System.Xml;
using Microsoft.AspNetCore.Components.Web;

namespace ListenBucket.Services.Podcast;

public class PodcastSyndicationFeed
{
    private const string iTunesNamespace = "http://www.itunes.com/dtds/podcast-1.0.dtd";
    public required string Title { get; set; }
    public required string Author { get; set; }
    public string? Description { get; set; }
    public Uri? AlternativeLink { get; set; }
    public DateTimeOffset LastUpdatedTime { get; set; }
    public required IEnumerable<PodcastCategory> Categories { get; set; }
    public string? Language { get; set; }
    public bool? Explicit { get; set; }
    public PodcastType PodcastType { get; set; } = PodcastType.Episodic;
    public Uri? ImageUrl { get; set; }
    public List<PodcastItem>? Items { get; set; }

    public async Task WriteToStreamAsync(Stream destination)
    {
        var settings = new XmlWriterSettings
        {
            Indent = true,
            OmitXmlDeclaration = false,
            Async = true
        };
        
        await using var writer = XmlWriter.Create(destination, settings);
        // Write the <rss> root with namespaces
        await writer.WriteStartDocumentAsync();
        writer.WriteStartElement("rss");
        writer.WriteAttributeString("version", "2.0");
        await writer.WriteAttributeStringAsync("xmlns", "itunes", null, iTunesNamespace);
        await writer.WriteAttributeStringAsync("xmlns", "podcast", null, "https://podcastindex.org/namespace/1.0");
        await writer.WriteAttributeStringAsync("xmlns", "atom", null, "http://www.w3.org/2005/Atom");
        await writer.WriteAttributeStringAsync("xmlns", "media", null, "http://search.yahoo.com/mrss/");
        await writer.WriteAttributeStringAsync("xmlns", "googleplay", null, "http://www.google.com/schemas/play-podcasts/1.0");
        await writer.WriteAttributeStringAsync("xmlns", "content", null, "http://purl.org/rss/1.0/modules/content/");
        writer.WriteStartElement("channel");
        
        writer.WriteStartElement("title");
        await writer.WriteStringAsync(Title);
        await writer.WriteEndElementAsync(); // </title>

        await writer.WriteNodeIfNotNull("language", Language);
        await writer.WriteNodeIfNotNull("description", Description);
        
        if (Explicit != null)
        {
            await writer.WriteStartElementAsync(prefix: "itunes", "explicit", iTunesNamespace);
            await writer.WriteStringAsync(Explicit == true ? "yes" : "no");
            await writer.WriteEndElementAsync();
        }
        
        await writer.WriteStartElementAsync(prefix: "itunes", "type", iTunesNamespace);
        await writer.WriteStringAsync(PodcastType.ToString().ToLower());
        await writer.WriteEndElementAsync();
        
        await writer.WriteStartElementAsync(prefix: "itunes", "author", iTunesNamespace);
        await writer.WriteStringAsync(Author.ToLower());
        await writer.WriteEndElementAsync();

        foreach (var category in Categories)
        {
            await writer.WriteStartElementAsync(prefix: "itunes", "category", iTunesNamespace);
            await writer.WriteAttributeStringAsync(prefix: null, localName: "text", value: category.ToString(), ns: null); 
            await writer.WriteEndElementAsync(); // </itunes:category>
        }

        if (ImageUrl != null)
        {
            writer.WriteStartElement("image");
            
            writer.WriteStartElement("url");
            await writer.WriteStringAsync(ImageUrl.ToString());
            await writer.WriteEndElementAsync();// </url>
            
            writer.WriteStartElement("title");
            await writer.WriteStringAsync(Title);
            await writer.WriteEndElementAsync(); // </title>
            
            await writer.WriteEndElementAsync(); // </image>
        }

        if (Items != null)
        {
            foreach (var item in Items)
            {
                writer.WriteStartElement("item");
                
                writer.WriteStartElement("guid");
                await writer.WriteAttributeStringAsync(prefix: null, localName: "isPermaLink", ns: null, value: "false");
                await writer.WriteCDataAsync(item.Id);
                await writer.WriteEndElementAsync();
                
                await writer.WriteNodeIfNotNull("title", item.Title);
                await writer.WriteNodeIfNotNull("description", item.Description);
                await writer.WriteNodeIfNotNull("pubDate", item.PublishedDate);

                if (item.ImageUrl != null)
                {
                    await writer.WriteStartElementAsync(prefix: "itunes", localName: "image", ns: iTunesNamespace);
                    writer.WriteAttributeString("href", item.ImageUrl.ToString());
                    await writer.WriteEndElementAsync(); // </itunes:image>
                }
                
                writer.WriteStartElement("enclosure");
                writer.WriteAttributeString("url", item.MediaUrl.ToString());
                writer.WriteAttributeString("type", "audio/mpeg");
                // Length is optional, set to 0 if unknown
                writer.WriteAttributeString("length", "0");
                await writer.WriteEndElementAsync(); // </enclosure>

                await writer.WriteNodeIfNotNull(prefix: "itunes", localName: "summary", ns: iTunesNamespace, value: item.Description);
                
                await writer.WriteEndElementAsync();// </item>
            }
        }
        
        await writer.WriteEndElementAsync(); // </channel>
        await writer.WriteEndElementAsync(); // </rss>
        await writer.WriteEndDocumentAsync();
    }
}

public class PodcastItem
{
    public required string Title { get; set; }
    public string? Description { get; set; }
    public DateTimeOffset? PublishedDate { get; set; }
    public DateTimeOffset? LastUpdatedTime { get; set; }
    public required string Id { get; set; }
    public required Uri MediaUrl { get; set; }
    public Uri? ImageUrl { get; set; }
}

public enum PodcastType
{
    Episodic,
    Serial
}

public enum PodcastCategory
{
    Arts,
    Business,
    Comedy,
    Education,
    Fiction,
    Government,
    HealthAndFitness,
    History,
    KidsAndFamily,
    Leisure,
    Music,
    NewsAndPolitics,
    ReligionAndSpirituality,
    Science,
    SocietyAndCulture,
    Sports,
    Technology,
    TVAndFilm
}

public static class XmlWriterExtensions
{
    public static async Task WriteNodeIfNotNull(this XmlWriter writer, string? prefix, string localName, string? ns, string? value)
    {
        if (!string.IsNullOrWhiteSpace(value))
        {
            await writer.WriteStartElementAsync(prefix, localName, ns);
            await writer.WriteStringAsync(value);
            await writer.WriteEndElementAsync();
        }
    }
    
    public static async Task WriteNodeIfNotNull(this XmlWriter writer, string name, string? value)
    {
        if (!string.IsNullOrWhiteSpace(value))
        {
            writer.WriteStartElement(name);
            await writer.WriteStringAsync(value);
            await writer.WriteEndElementAsync();
        }
    }
    
    public static async Task WriteNodeIfNotNull(this XmlWriter writer, string name, DateTimeOffset? value)
    {
        if (value != null)
        {
            writer.WriteStartElement(name);
            await writer.WriteStringAsync(value.Value.ToString("R"));
            await writer.WriteEndElementAsync();
        }
    }
}