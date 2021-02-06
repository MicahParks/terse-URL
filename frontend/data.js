function Terse(javascriptTracking, mediaPreview, originalURL, redirectType, shortenedURL) {
    this.javascriptTracking = javascriptTracking;
    this.mediaPreview = mediaPreview;
    this.originalURL = originalURL;
    this.redirectType = redirectType;
    this.shortenedURL = shortenedURL;
}

function MediaPreview(inherit, og, title, twitter) {
    this.inherit = inherit
    this.title = title;
    this.twitter = twitter;
    this.og = og;
}

function TerseSummary(originalURL, shortenedURL, redirectType, visitCount) {
    this.originalURL = originalURL;
    this.redirectType = redirectType;
    this.shortenedURL = shortenedURL;
    this.visitCount = visitCount;
}