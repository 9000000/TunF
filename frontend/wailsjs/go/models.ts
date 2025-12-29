export namespace main {
	
	export class Config {
	    lastListenPort: string;
	    lastTargetAddr: string;
	    autoOpenFirewall: boolean;
	    autoStart: boolean;
	    history: string[];
	    targetHistory: string[];
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.lastListenPort = source["lastListenPort"];
	        this.lastTargetAddr = source["lastTargetAddr"];
	        this.autoOpenFirewall = source["autoOpenFirewall"];
	        this.autoStart = source["autoStart"];
	        this.history = source["history"];
	        this.targetHistory = source["targetHistory"];
	    }
	}

}

