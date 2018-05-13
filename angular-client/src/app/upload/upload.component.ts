import { Component, OnInit, ElementRef, ViewChild } from '@angular/core';
import { HttpClient, HttpResponse, HttpEventType } from '@angular/common/http';
import { FormGroup, FormBuilder, FormControl } from '@angular/forms';
import { PercentPipe } from '@angular/common';
import { FileUploadService } from '../services/file-upload.service';
import { AfterViewInit } from '@angular/core/src/metadata/lifecycle_hooks';

import * as _ from 'lodash';

// Model
import { Post, Label } from '../model/post';

@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.css']
})
export class UploadComponent implements AfterViewInit {

  selectedFiles: FileList;
  currentFileUpload: File;
  progress: { percentage: number } = { percentage: 0 };
  imageSrc: any;
  fileName: string;
  post: Post;
  label: Label;
  openModal = false;

 mockLabel: Label = {
   label: 'Bear',
   probability: 0.9834529
 };

  mockData: Post = {
    filename: 'bear.jpg',
    labels: [this.mockLabel, this.mockLabel, this.mockLabel]
  };

  @ViewChild('fileUpload') fileUploader: ElementRef;

  constructor(private uploadService: FileUploadService) {}

  ngAfterViewInit() {
    console.log(this.fileUploader.nativeElement.value);
  }

  displayPhoto(fileInput) {
    if (fileInput.target.files && fileInput.target.files[0]) {
      const reader = new FileReader();

      reader.onload = ((e) => {
        this.imageSrc = e.target['result'];
      });

      reader.readAsDataURL(fileInput.target.files[0]);
    }
  }

  // onChange(event) {
  //   this.selectedFiles = event.target.files;
  //   console.log(this.selectedFiles);
  //   const reader = new FileReader();
  //   reader.onload = (e: any) => {
  //     this.logo = e.target.result;
  //   };
  //   reader.readAsDataURL(event.target.files[0]);
  //   this.fileName = event.target.files[0].name;
  //   this.fileName = this.selectedFiles.item(0).name;
  // }


  onChange(event) {
    this.selectedFiles = event.target.files;

    this.fileName = this.selectedFiles.item(0).name;
  }

  clearAlertMessage() {
    this.post = null;
    this.fileUploader.nativeElement.click();
    this.fileUploader.nativeElement.value = '';
  }

  clear() {
    this.selectedFiles = undefined;
    this.fileName = '';
    this.fileUploader.nativeElement.value = '';
  }


  upload() {
    this.progress.percentage = 0;

    this.currentFileUpload = this.selectedFiles.item(0);

    this.uploadService.uploadFile(this.currentFileUpload).subscribe(event => {
      if (event.type === HttpEventType.UploadProgress) {
        this.progress.percentage = Math.round(100 * event.loaded / event.total);
        console.log(this.progress.percentage);
      } else if ( event instanceof HttpResponse ) {
        console.log(event.body);
      }
    });

    this.selectedFiles = undefined;
    this.fileName = '';
    this.fileUploader.nativeElement.value = '';
  }

  analyze() {
    this.currentFileUpload = this.selectedFiles.item(0);

    this.uploadService.analyzePhoto(this.currentFileUpload).subscribe(data => {
      console.log(data);
      this.post = data;
      this.openModal = !this.openModal;
    });

    this.selectedFiles = undefined;
    this.fileName = '';
    this.fileUploader.nativeElement.value = '';
  }

  openDialog() {
    this.openModal = !this.openModal;
    this.post = this.mockData;

    this.selectedFiles = undefined;
    this.fileName = '';
    this.fileUploader.nativeElement.value = '';
  }
}
