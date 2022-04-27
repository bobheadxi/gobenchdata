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

function FromArray<T>(Ctor: { new (v: any): T }, data?: any[] | any, def = null): T[] | null {
	if (!data || !Object.keys(data).length) return def;
	const d = Array.isArray(data) ? data : [data];
	return d.map((v: any) => new Ctor(v));
}

function ToObject(o: any, typeOrCfg: any = {}, child = false): any {
	if (o == null) return null;
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
		if (v === undefined) continue;
		if (v === null) continue;
		d[k] = ToObject(v, typeOrCfg[k] || {}, true);
	}

	return d;
}

// structs
// struct2ts:go.bobheadxi.dev/gobenchdata/web.ChartDisplay
class ChartDisplay {
  fullWidth: boolean;

  constructor(data?: any) {
    const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
    this.fullWidth = ('fullWidth' in d) ? d.fullWidth as boolean : false;
  }

  toObject(): any {
    const cfg: any = {};
    return ToObject(this, cfg);
  }
}

// struct2ts:go.bobheadxi.dev/gobenchdata/web.Chart
class Chart {
  name: string;
  description: string;
  package: string;
  benchmarks: string[] | null;
  metrics: { [key: string]: boolean };
  display: ChartDisplay | null;

  constructor(data?: any) {
    const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
    this.name = ('name' in d) ? d.name as string : '';
    this.description = ('description' in d) ? d.description as string : '';
    this.package = ('package' in d) ? d.package as string : '';
    this.benchmarks = ('benchmarks' in d) ? d.benchmarks as string[] : null;
    this.metrics = ('metrics' in d) ? d.metrics as { [key: string]: boolean } : {};
    this.display = ('display' in d) ? new ChartDisplay(d.display) : null;
  }

  toObject(): any {
    const cfg: any = {};
    return ToObject(this, cfg);
  }
}

// struct2ts:go.bobheadxi.dev/gobenchdata/web.ChartGroup
class ChartGroup {
  name: string;
  description: string;
  charts: Chart[] | null;

  constructor(data?: any) {
    const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
    this.name = ('name' in d) ? d.name as string : '';
    this.description = ('description' in d) ? d.description as string : '';
    this.charts = Array.isArray(d.charts) ? d.charts.map((v: any) => new Chart(v)) : null;
  }

  toObject(): any {
    const cfg: any = {};
    return ToObject(this, cfg);
  }
}

// struct2ts:go.bobheadxi.dev/gobenchdata/web.Config
class Config {
  title: string;
  description: string;
  repository: string;
  benchmarksFile: string | null;
  chartGroups: ChartGroup[] | null;

  constructor(data?: any) {
    const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
    this.title = ('title' in d) ? d.title as string : '';
    this.description = ('description' in d) ? d.description as string : '';
    this.repository = ('repository' in d) ? d.repository as string : '';
    this.benchmarksFile = ('benchmarksFile' in d) ? d.benchmarksFile as string : null;
    this.chartGroups = Array.isArray(d.chartGroups) ? d.chartGroups.map((v: any) => new ChartGroup(v)) : null;
  }

  toObject(): any {
    const cfg: any = {};
    return ToObject(this, cfg);
  }
}

// struct2ts:go.bobheadxi.dev/gobenchdata/bench.Mem
class Mem {
  BytesPerOp: number;
  AllocsPerOp: number;
  MBPerSec: number;

  constructor(data?: any) {
    const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
    this.BytesPerOp = ('BytesPerOp' in d) ? d.BytesPerOp as number : 0;
    this.AllocsPerOp = ('AllocsPerOp' in d) ? d.AllocsPerOp as number : 0;
    this.MBPerSec = ('MBPerSec' in d) ? d.MBPerSec as number : 0;
  }

  toObject(): any {
    const cfg: any = {};
    cfg.BytesPerOp = 'number';
    cfg.AllocsPerOp = 'number';
    cfg.MBPerSec = 'number';
    return ToObject(this, cfg);
  }
}

// struct2ts:go.bobheadxi.dev/gobenchdata/bench.Benchmark
class Benchmark {
  Name: string;
  Runs: number;
  NsPerOp: number;
  Mem: Mem;
  Custom?: { [key: string]: number };

  constructor(data?: any) {
    const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
    this.Name = ('Name' in d) ? d.Name as string : '';
    this.Runs = ('Runs' in d) ? d.Runs as number : 0;
    this.NsPerOp = ('NsPerOp' in d) ? d.NsPerOp as number : 0;
    this.Mem = new Mem(d.Mem);
    this.Custom = ('Custom' in d) ? d.Custom as { [key: string]: number } : {};
  }

  toObject(): any {
    const cfg: any = {};
    cfg.Runs = 'number';
    cfg.NsPerOp = 'number';
    return ToObject(this, cfg);
  }
}

// struct2ts:go.bobheadxi.dev/gobenchdata/bench.Suite
class Suite {
  Goos: string;
  Goarch: string;
  Pkg: string;
  Benchmarks: Benchmark[] | null;

  constructor(data?: any) {
    const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
    this.Goos = ('Goos' in d) ? d.Goos as string : '';
    this.Goarch = ('Goarch' in d) ? d.Goarch as string : '';
    this.Pkg = ('Pkg' in d) ? d.Pkg as string : '';
    this.Benchmarks = Array.isArray(d.Benchmarks) ? d.Benchmarks.map((v: any) => new Benchmark(v)) : null;
  }

  toObject(): any {
    const cfg: any = {};
    return ToObject(this, cfg);
  }
}

// struct2ts:go.bobheadxi.dev/gobenchdata/bench.Run
class Run {
  Version?: string;
  Date: number;
  Tags?: string[] | null;
  Suites: Suite[] | null;

  constructor(data?: any) {
    const d: any = (data && typeof data === 'object') ? ToObject(data) : {};
    this.Version = ('Version' in d) ? d.Version as string : '';
    this.Date = ('Date' in d) ? d.Date as number : 0;
    this.Tags = ('Tags' in d) ? d.Tags as string[] : null;
    this.Suites = Array.isArray(d.Suites) ? d.Suites.map((v: any) => new Suite(v)) : null;
  }

  toObject(): any {
    const cfg: any = {};
    cfg.Date = 'number';
    return ToObject(this, cfg);
  }
}

// exports
export {
  ChartDisplay,
  Chart,
  ChartGroup,
  Config,
  Mem,
  Benchmark,
  Suite,
  Run,
  ParseDate,
  ParseNumber,
  FromArray,
  ToObject,
};
