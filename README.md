boltcopy
========

`boltcopy` is a simple tool for copying [BoltDB](https://github.com/boltdb/bolt) database files optionally white listing or black listing buckets. If you wanted to reclaim used space or easily remove buckets from databases en bulk this might be the tool for you.

**Note:** `boltcopy` does not support nested buckets

**Note:** `boltcopy` assumes your bucket names are human-readable strings

Usage
-----

```sh
$ ./boltcopy -h
Usage:
./boltcopy [flags] <intput.db> <output.db>
Flags:
  -b value
    	name of bucket (can specify multiple)
  -i	include provided buckets (i.e. white list behavior)
```

By default `boltcopy` will simply copy all key/value pairs from the input database to the output database. However you can optionally include or exclude certain buckets by name. By default we operate in "black list" mode where any specified bucket names will not be copied to the destination database.

For example if the BoltDB `in.db` has three buckets: "A", "B", and "C".

Executing:

```sh
boltcopy -b A -b B in.db out.db
```

The buckets that end up in `out.db` would just be "C". Alternatively if we to run in "white list" mode we can invert our call with the `-i` switch:

```sh
boltcopy -i -b A -b B in.db out.db
```

In this case only the "A" and "B" buckets will be copied. Buckets that don't match or are missing in either mode will be ignored.