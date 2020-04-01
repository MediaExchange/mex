export enum MediaType {
    Movie,
    TvShow,
}

export function toString(t: MediaType) {
    switch(t) {
        case MediaType.Movie:
            return 'Movie';
        case MediaType.TvShow:
            return 'TvShow';
    }
}
