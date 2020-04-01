import { Component, OnDestroy, OnInit } from '@angular/core';
import {ActivatedRoute, NavigationEnd, Router, RouterEvent} from '@angular/router';
import { Subscription } from 'rxjs';
import { SearchResult } from '../../models/search-result';
import { SearchService } from '../../services/search.service';
import {DetailsService} from "../../services/details.service";
import {DetailsResult} from "../../models/details-result";

@Component({
  selector: 'app-search-result',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.less']
})
export class SearchComponent implements OnDestroy, OnInit {
  private refresh: Subscription;
  public detailsModalVisible: boolean;
  public results: Array<SearchResult>;
  public details: DetailsResult;

  // DI Constructor
  constructor(private activatedRoute: ActivatedRoute,
              private router: Router,
              private searchService: SearchService,
              private detailsService: DetailsService) {
  }

  ngOnDestroy() {
    // Have to unsubscribe to prevent memory leaks.
    if (this.refresh) {
      this.refresh.unsubscribe();
    }
  }

  ngOnInit(): void {
    // This allows the search page to route to itself with a new query string.
    this.router.onSameUrlNavigation = 'reload';
    this.refresh = this.router.events.subscribe((e: any) => {
      if (e instanceof NavigationEnd) {
        this.getSearchResults();
      }
    });

    this.details = new DetailsResult();
    this.getSearchResults();
  }

  public getSearchResults() {
    const q = this.activatedRoute.snapshot.queryParamMap.get('q');
    if (q !== undefined && q !== null && q !== '') {
      this.searchService.search(q)
          .subscribe(
              (data: SearchResult[]) => {
                console.log(data);
                this.results = data;
              },
              (err: any) => {
                console.error(err);
              });
    }
  }

  public showDetails(sr: SearchResult) {
    this.detailsService.details(sr.id)
        .subscribe(
            (data: DetailsResult) => {
              console.log(data);
              this.details = data;
              this.detailsModalVisible = true;
            },
            (err: any) => {
              console.error(err);
            });
  }

  public sortedResults() {
    return this.results.sort((l: SearchResult, r: SearchResult) => new Date(r.releaseDate).valueOf() - new Date(l.releaseDate).valueOf());
  }
}
