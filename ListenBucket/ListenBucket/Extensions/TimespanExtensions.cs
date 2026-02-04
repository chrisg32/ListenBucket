namespace ListenBucket.Extensions;

public static class TimespanExtensions
{
    public static string Format(this TimeSpan? timeSpan)
    {
        if (timeSpan == null) return string.Empty;

        if (timeSpan.Value.Hours >= 1)
        {
            return timeSpan.Value.Minutes > 0 ? $"{timeSpan.Value.Hours} hr {timeSpan.Value.Minutes} min" : $"{timeSpan.Value.Hours} hr";
        }
        return $"{timeSpan.Value.Minutes} min";
    }
}