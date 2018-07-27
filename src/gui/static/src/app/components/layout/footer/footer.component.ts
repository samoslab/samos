import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { MatDialog, MatDialogConfig } from '@angular/material/dialog';
import { SendSamosComponent } from '../../pages/send-samos/send-samos.component';

@Component({
  selector: 'app-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss']
})
export class FooterComponent {

  constructor(
    private dialog: MatDialog,
    private router: Router,
  ) { }

  openWalletsPage() {
    this.router.navigate(['/wallets']);
  }

  openSendPage() {
    const config = new MatDialogConfig();
    config.width = '566px';
    this.dialog.open(SendSamosComponent, config);
  }

  openTransactions() {
    this.router.navigate(['/transactions']);
  }
}
