export namespace folderselector {
	
	export class FolderSelectorResult {
	    Directory: string;
	    Files: string[];
	
	    static createFrom(source: any = {}) {
	        return new FolderSelectorResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Directory = source["Directory"];
	        this.Files = source["Files"];
	    }
	}

}

