import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { DownloadsComponent } from './components/downloads/downloads.component';
import { SearchComponent } from './components/search/search.component';

const routes: Routes = [
  {
    path: 'downloads',
    component: DownloadsComponent
  },
  {
    path: 'search',
    component: SearchComponent
  },
  {
    path: '',
    redirectTo: '/downloads',
    pathMatch: 'full'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
