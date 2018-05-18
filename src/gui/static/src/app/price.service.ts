import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

@Injectable()
export class PriceService {

  price: Subject<number> = new BehaviorSubject<number>(null);

  constructor(
    private http: Http,
  ) {
           setTimeout( ()=> {
                this.price.next(0)
              }, 0);
  }

}
