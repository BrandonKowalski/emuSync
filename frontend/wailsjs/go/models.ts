export namespace models {
	
	export class Device {
	    id: string;
	    adb_index: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Device(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.adb_index = source["adb_index"];
	        this.name = source["name"];
	    }
	}

}

