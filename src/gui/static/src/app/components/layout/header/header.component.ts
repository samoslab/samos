import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { PriceService } from '../../../price.service';
import { Subscription } from 'rxjs/Subscription';
import { WalletService } from '../../../services/wallet.service';
import { BlockchainService } from '../../../services/blockchain.service';
import { Observable } from 'rxjs/Observable';
import { ApiService } from '../../../services/api.service';
import { Http } from '@angular/http';
import { AppService } from '../../../services/app.service';
import { IntervalObservable } from 'rxjs/observable/IntervalObservable';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit, OnDestroy {
  @Input() title: string;
  @Input() coins: number;
  @Input() hours: number;

  current: number;
  highest: number;
  percentage: number;
  querying = true;
  version: string;
  releaseVersion: string;
  updateAvailable: boolean;
  hasPendingTxs: boolean;

  private price: number;
  private priceSubscription: Subscription;
  private walletSubscription: Subscription;

  get balance() {
    if (this.price === null) { return 'loading..'; }
    const balance = Math.round(this.coins * this.price * 100) / 100;
    return '$' + balance.toFixed(2) + ' ($' + (Math.round(this.price * 100) / 100) + ')';
  }

  get loading() {
    return !this.current || !this.highest || this.current !== this.highest;
  }

  constructor(
    public appService: AppService,
    private apiService: ApiService,
    private blockchainService: BlockchainService,
    private priceService: PriceService,
    private walletService: WalletService,
    private http: Http,
  ) { }

  ngOnInit() {
    this.setVersion();
    this.priceSubscription = this.priceService.price.subscribe(price => this.price = price);
    this.walletSubscription = this.walletService.allAddresses().subscribe(addresses => {
      addresses = addresses.reduce((array, item) => {
        if (!array.find(addr => addr.address === item.address)) {
          array.push(item);
        }
        return array;
      }, []);

      this.coins = addresses.map(addr => addr.coins >= 0 ? addr.coins : 0).reduce((a, b) => a + b, 0);
      this.hours = addresses.map(addr => addr.hours >= 0 ? addr.hours : 0).reduce((a, b) => a + b, 0);
    });

    this.blockchainService.progress
      .filter(response => !!response)
      .subscribe(response => {
        this.querying = false;
        this.highest = response.highest;
        this.current = response.current;
        this.percentage = this.current && this.highest ? (this.current / this.highest) : 0;
      });

    this.walletService.pendingTransactions().subscribe(txs => {
      this.hasPendingTxs = txs.length > 0;
    });
  }

  ngOnDestroy() {
    this.priceSubscription && this.priceSubscription.unsubscribe();
    this.walletSubscription && this.walletSubscription.unsubscribe();
  }

  setVersion() {
    // Set build version
    this.apiService.getVersion().first()
      .subscribe(output =>  {
        this.version = output.version;
        this.retrieveReleaseVersion();
      });
  }

  private toNum(a){
    var a=a.toString();
    var c=a.split('.');
    var num_place=["","0","00","000","0000"],r=num_place.reverse();
    for (var i=0;i<c.length;i++){
      var len=c[i].length;
      c[i]=r[len]+c[i];
    }
    var res= c.join('');
    return res;
  } 

  private higherVersion(first: string, second: string): boolean {
    console.log(first+' '+second);
    // const fa = first.split('.');
    // const fb = second.split('.');
    // for (let i = 0; i < 3; i++) {
    //   const na = Number(fa[i]);
    //   const nb = Number(fb[i]);
    //   if (na > nb || !isNaN(na) && isNaN(nb)) {
    //     return true;
    //   } else if (na < nb || isNaN(na) && !isNaN(nb)) {
    //     return false;
    //   }
    // }
    // return false;

    var _a=this.toNum(first),_b=this.toNum(second);
    if(_a==_b) return false;
    if(_a>_b) return true;  
    if(_a<_b) return false;  
  }

  private retrieveReleaseVersion() {
    //this.http.get('http://samos.io/version/samos-tags.txt')
    //this.http.get('https://api.github.com/repos/samoslab/samos/tags')
    // this.http.get('http://samos.io/api/version?tag=samos')
    //       .map((res: any) => res.json())
    //       .map((res: any) => {
    //         let r = res.json();
    //         return r.length < 1 ? [{name: ""}] : r;
    //       })
    //   .catch((error: any) => Observable.throw(error || 'Unable to fetch latest release version from github.'))
    //   .subscribe(response =>  {
    //     this.releaseVersion = response.find(element => element['name'].indexOf('rc') === -1)['name'].substr(1);
    //     this.updateAvailable = this.higherVersion(this.releaseVersion, this.version);
    //   });

    var _this = this;
    try{
      this.http.get('http://samos.io/api/version?tag=samos').map((res: any) => res.json()).subscribe(function (response) {
        response.version_code ? response : 'return';
        _this.releaseVersion = response.version_code;
        var releaseVersion = _this.releaseVersion.toString();
        _this.updateAvailable = _this.higherVersion(releaseVersion,_this.version);
  
      })
    }catch(error){
      Observable.throw(error || 'Unable to fetch latest release version from github.');
    }
  }
}
