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

