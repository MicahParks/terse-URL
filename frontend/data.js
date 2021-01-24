function Terse(javascriptTracking, mediaPreview, original, redirectType, shortened) {
    this.javascriptTracking = javascriptTracking;
    this.mediaPreview = mediaPreview;
    this.originalURL = original;
    this.redirectType = redirectType;
    this.shortenedURL = shortened;
}

function MediaPreview(og, title, twitter) {
    this.title = title;
    this.twitter = twitter;
    this.og = og;
}
