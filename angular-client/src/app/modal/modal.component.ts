import { Component, OnInit, OnChanges, Input, Output, EventEmitter } from '@angular/core';
import { trigger, state, style, animate, transition } from '@angular/animations';

import { Post, Label } from '../model/post';

@Component({
  selector: 'app-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.css'],
  animations: [
    trigger('modal',[
      transition('void => *', [
        style({ transform: 'scale3d(.3, .3, .3)'}),
        animate(100)
      ]),
      transition('* => void', [
        animate(100, style({ transform: 'scale3d(.0, .0, .0)'}))
      ])
    ])
  ]
})
export class ModalComponent implements OnInit {

  @Input() closable = true;
  @Input() visible: boolean;
  @Output() visibleChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  post: Post;

  constructor() { }

  ngOnInit() {
  }

  close() {
    this.visible = false;
    this.visibleChange.emit(this.visible);
  }




}
