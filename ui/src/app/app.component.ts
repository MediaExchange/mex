import { Component } from '@angular/core';
import { NgForm } from '@angular/forms';
import { Router } from '@angular/router';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.less']
})

export class AppComponent {
  title = 'MEX';

  // DI constructor.
  constructor(private router: Router) {
  }

  // Search form handler.
  onSearch(f: NgForm) {
    this.router.navigate(['search'], {queryParams: {q: f.value.name}});
  }
}
