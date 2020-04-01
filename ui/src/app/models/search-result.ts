import {MediaType} from './media-type.enum';

export class SearchResult {
    // ID of the media found.
    id: string;

    // Name of the media found.
    title: string;

    // When the media first aired on TV or was released in theaters.
    releaseDate: string;

    // URI of an image that can be displayed.
    posterUri: string;

    // Overview description of the media.
    overview: string;

    // Type of media.
    type: MediaType;

    // Whether the media is for adults
    adult: boolean;

    // constructor accepts an object and copies the fields into the new SearchResults instance.
    // This is used by `SearchService` to convert the JSON response to the actual object so the methods work
    // as expected.
    constructor(obj?: any) {
        Object.assign(this, obj);
    }

    // getShortOverview returns the media overview, which can be long, but truncates it to 130 characters.
    getShortOverview() {
        if (this.overview.length > 130) {
            return this.overview.substr(0, 130) + '...';
        }
        return this.overview;
    }

    // getTitle returns the title of the media with the year appeneded in parenthesis, e.g. "title (year)"
    public getTitle() {
        let t = this.title;
        if (this.hasReleaseDate()) {
            t += ' (' + this.getYear() + ')';
        }
        return t;
    }

    // getYear returns the year portion of the releaseData field.
    public getYear() {
        return new Date(this.releaseDate).getFullYear();
    }

    // hasReleaseDate verifies that the release date exists.
    public hasReleaseDate() {
        return this.releaseDate !== undefined && this.releaseDate != null && this.releaseDate.length > 0;
    }
}
