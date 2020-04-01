import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders, HttpParams} from '@angular/common/http';
import {Observable} from 'rxjs';
import {SearchResult} from '../models/search-result';
import {map} from 'rxjs/operators';
import {DetailsResult} from '../models/details-result';

@Injectable({
    providedIn: 'root'
})
export class DetailsService {
    private detailsUrl = 'http://localhost:9000/api/details';

    constructor(private http: HttpClient) {
    }

    details(id: string): Observable<DetailsResult> {
        const headers = new HttpHeaders()
            .append('Accept', 'application/json');
        const params = new HttpParams()
            .append('id', id);
        return this.http.get<DetailsResult>(this.detailsUrl, {headers, params}).pipe(
            map(r => new DetailsResult(r))
        );
    }
}
