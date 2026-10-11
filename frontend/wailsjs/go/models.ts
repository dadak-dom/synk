export namespace folderselector {
	
	export class FolderSelectorResult {
	    Directory: string;
	    Files: string[];
	    Folders: string[];
	
	    static createFrom(source: any = {}) {
	        return new FolderSelectorResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Directory = source["Directory"];
	        this.Files = source["Files"];
	        this.Folders = source["Folders"];
	    }
	}

}

export namespace network {
	
	export class RemotePeer {
	    IP: string;
	    DeviceName: string;
	
	    static createFrom(source: any = {}) {
	        return new RemotePeer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.IP = source["IP"];
	        this.DeviceName = source["DeviceName"];
	    }
	}

}

