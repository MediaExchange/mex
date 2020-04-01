import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders, HttpParams} from '@angular/common/http';
import {Observable} from 'rxjs';
import {SearchResult} from '../models/search-result';
import {map} from 'rxjs/operators';

@Injectable({
    providedIn: 'root'
})
export class SearchService {
    private searchUrl = 'http://localhost:9000/api/search';

    constructor(private http: HttpClient) {
    }

    search(name: string): Observable<Array<SearchResult>> {
        const headers = new HttpHeaders()
            .append('Accept', 'application/json');
        const params = new HttpParams()
            .append('q', name);
        return this.http.get<Array<SearchResult>>(this.searchUrl, {headers, params}).pipe(
            map(a => a.map(r => new SearchResult(r)))
        );
    }
}
