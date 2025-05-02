export namespace downloader {
	
	export class DownloadTask {
	    id: string;
	    url: string;
	    name: string;
	    status: string;
	    savePath: string;
	    // Go type: time
	    startTime: any;
	    // Go type: time
	    completeTime: any;
	    error: string;
	    // Go type: struct { Current int "json:\"current\""; Total int "json:\"total\"" }
	    progress: any;
	
	    static createFrom(source: any = {}) {
	        return new DownloadTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.url = source["url"];
	        this.name = source["name"];
	        this.status = source["status"];
	        this.savePath = source["savePath"];
	        this.startTime = this.convertValues(source["startTime"], null);
	        this.completeTime = this.convertValues(source["completeTime"], null);
	        this.error = source["error"];
	        this.progress = this.convertValues(source["progress"], Object);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace getter {
	
	export class Manga {
	    name: string;
	    path: string;
	    previewImg: string;
	    imagesCount: number;
	    images?: string[];
	
	    static createFrom(source: any = {}) {
	        return new Manga(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.previewImg = source["previewImg"];
	        this.imagesCount = source["imagesCount"];
	        this.images = source["images"];
	    }
	}

}

export namespace storage {
	
	export class Storage {
	
	
	    static createFrom(source: any = {}) {
	        return new Storage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

