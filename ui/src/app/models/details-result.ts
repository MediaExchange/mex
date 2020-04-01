import { MediaType } from './media-type.enum';

// Episode contains detailed information about an episode within a series.
export class Episode {
    // Name of the episode.
    name: string;

    // Absolute episode number across all seasons.
    number: number;

    // Season number.
    season: number;

    // Original date the episode aired.
    airDate: string;

    // Episode number within the season.
    episode: number;

    // Overview description of the episode.
    overview: string;
}

// Link contains a reference to an external source of information about the media.
export class Link {
    // Name of the external source.
    name: string;

    // Link to the media reference.
    url: string;
}

export class DetailsResult {
    // ID of the media.
    id: string;

    // Type of media (Movie, TV Show, etc.)
    type: MediaType;

    // Whether the media is for adults (Rated X, TV-MA, etc.)
    adult: boolean;

    // Links to external information about the media.
    links: Array<Link>;

    // Title of the media.
    title: string;

    // Current status of the media (Released, In Production, etc.)
    status: string;

    // Media runtime in minutes.
    runtime: number;

    // Detailed information about each episode in a series.
    episodes: Array<Episode>;

    // Overview description of the media.
    overview: string;

    // URL of the media poster.
    posterUri: string;

    // Date the media was first seen in theaters or on TV.
    releaseDate: string;

    // constructor accepts an object and copies the fields into the new SearchResults instance.
    // This is used by `SearchService` to convert the JSON response to the actual object so the methods work
    // as expected.
    constructor(obj?: any) {
        Object.assign(this, obj);
    }

    // getEpisodes returns all of the episodes within a season.
    public getEpisodes(season: number): Array<Episode> {
        return this.episodes.filter(e => e.season === season).sort((l, r) => l.episode - r.episode);
    }

    // getSeasons returns an array of the season numbers.
    public getSeasons(): Array<number> {
        return Array.from(new Set(this.episodes.map(e => e.season))).sort((l, r) => l - r);
    }

    // getTitle returns the title of the media with the year appeneded in parenthesis, e.g. "title (year)"
    public getTitle(): string {
        let t = this.title;
        if (this.hasReleaseDate()) {
            t += ' (' + this.getYear() + ')';
        }
        return t;
    }

    // getYear returns the year portion of the releaseData field.
    public getYear(): number {
        return new Date(this.releaseDate).getFullYear();
    }

    // hasEpisodes checks whether the details contains episodes.
    public hasEpisodes(): boolean {
        return this.episodes !== undefined && this.episodes !== null && this.episodes.length > 0;
    }

    // hasLinks checks whether the details contains links to external sources.
    public hasLinks(): boolean {
        return this.links !== undefined && this.links !== null && this.links.length > 0;
    }

    // hasReleaseDate verifies that the release date exists.
    public hasReleaseDate(): boolean {
        return this.releaseDate !== undefined && this.releaseDate !== null && this.releaseDate.length > 0;
    }
}
