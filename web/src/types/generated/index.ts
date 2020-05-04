// this file was automatically generated, DO NOT EDIT

// helpers
const maxUnixTSInSeconds = 9999999999;

function ParseDate(d: Date | number | string): Date {
	if (d instanceof Date) return d;
	if (typeof d === 'number') {
		if (d > maxUnixTSInSeconds) return new Date(d);
		return new Date(d * 1000); // go ts
	}
	return new Date(d);
}

function ParseNumber(v: number | string, isInt = false): number {
	if (!v) return 0;
	if (typeof v === 'number') return v;
	return (isInt ? parseInt(v) : parseFloat(v)) || 0;
}

function FromArray<T>(Ctor: { new(v: any): T }, data?: any[] | any, def = null): T[] | null {
	if (!data || !Object.keys(data).length) return def;
	const d = Array.isArray(data) ? data : [data];
	return d.map((v: any) => new Ctor(v));
}

function ToObject(o: any, typeOrCfg: any = {}, child = false): any {
	if (!o) return null;
	if (typeof o.toObject === 'function' && child) return o.toObject();

	switch (typeof o) {
		case 'string':
			return typeOrCfg === 'number' ? ParseNumber(o) : o;
		case 'boolean':
		case 'number':
			return o;
	}

	if (o instanceof Date) {
		return typeOrCfg === 'string' ? o.toISOString() : Math.floor(o.getTime() / 1000);
	}

	if (Array.isArray(o)) return o.map((v: any) => ToObject(v, typeOrCfg, true));

	const d: any = {};

	for (const k of Object.keys(o)) {
		const v: any = o[k];
		if (!v) continue;
		d[k] = ToObject(v, typeOrCfg[k] || {}, true);
	}

	return d;
}

// classes
// struct2ts:go.bobheadxi.dev/gobenchdata/web.ConfigChartConfig
class ConfigChartConfig {
	Name: string | null;
	Description: string | null;
	Package: string;
	Benchmarks: string[];

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.Name = ('Name' in d) ? d.Name as string : null;
		this.Description = ('Description' in d) ? d.Description as string : null;
		this.Package = ('Package' in d) ? d.Package as string : '';
		this.Benchmarks = ('Benchmarks' in d) ? d.Benchmarks as string[] : [];
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:go.bobheadxi.dev/gobenchdata/web.Config
class Config {
	Title: string;
	Description: string;
	BenchmarksFile: string | null;
	Charts: ConfigChartConfig[];

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.Title = ('Title' in d) ? d.Title as string : '';
		this.Description = ('Description' in d) ? d.Description as string : '';
		this.BenchmarksFile = ('BenchmarksFile' in d) ? d.BenchmarksFile as string : null;
		this.Charts = Array.isArray(d.Charts) ? d.Charts.map((v: any) => new ConfigChartConfig(v)) : [];
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:go.bobheadxi.dev/gobenchdata/bench.RunSuiteBenchmarkMem
class RunSuiteBenchmarkMem {
	BytesPerOp: number;
	AllocsPerOp: number;

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.BytesPerOp = ('BytesPerOp' in d) ? d.BytesPerOp as number : 0;
		this.AllocsPerOp = ('AllocsPerOp' in d) ? d.AllocsPerOp as number : 0;
	}

	toObject(): any {
		const cfg: any = {};
		cfg.BytesPerOp = 'number';
		cfg.AllocsPerOp = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:go.bobheadxi.dev/gobenchdata/bench.RunSuiteBenchmark
class RunSuiteBenchmark {
	Name: string;
	Runs: number;
	NsPerOp: number;
	Mem: RunSuiteBenchmarkMem;
	Custom: { [key: string]: number };

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.Name = ('Name' in d) ? d.Name as string : '';
		this.Runs = ('Runs' in d) ? d.Runs as number : 0;
		this.NsPerOp = ('NsPerOp' in d) ? d.NsPerOp as number : 0;
		this.Mem = new RunSuiteBenchmarkMem(d.Mem);
		this.Custom = ('Custom' in d) ? d.Custom as { [key: string]: number } : {};
	}

	toObject(): any {
		const cfg: any = {};
		cfg.Runs = 'number';
		cfg.NsPerOp = 'number';
		return ToObject(this, cfg);
	}
}

// struct2ts:go.bobheadxi.dev/gobenchdata/bench.RunSuite
class RunSuite {
	Goos: string;
	Goarch: string;
	Pkg: string;
	Benchmarks: RunSuiteBenchmark[];

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.Goos = ('Goos' in d) ? d.Goos as string : '';
		this.Goarch = ('Goarch' in d) ? d.Goarch as string : '';
		this.Pkg = ('Pkg' in d) ? d.Pkg as string : '';
		this.Benchmarks = Array.isArray(d.Benchmarks) ? d.Benchmarks.map((v: any) => new RunSuiteBenchmark(v)) : [];
	}

	toObject(): any {
		const cfg: any = {};
		return ToObject(this, cfg);
	}
}

// struct2ts:go.bobheadxi.dev/gobenchdata/bench.Run
class Run {
	Version: string;
	Date: number;
	Tags: string[];
	Suites: RunSuite[];

	constructor(data?: any) {
		const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
		this.Version = ('Version' in d) ? d.Version as string : '';
		this.Date = ('Date' in d) ? d.Date as number : 0;
		this.Tags = ('Tags' in d) ? d.Tags as string[] : [];
		this.Suites = Array.isArray(d.Suites) ? d.Suites.map((v: any) => new RunSuite(v)) : [];
	}

	toObject(): any {
		const cfg: any = {};
		cfg.Date = 'number';
		return ToObject(this, cfg);
	}
}

// exports
export {
	ConfigChartConfig,
	Config,
	RunSuiteBenchmarkMem,
	RunSuiteBenchmark,
	RunSuite,
	Run,
	ParseDate,
	ParseNumber,
	FromArray,
	ToObject,
};
