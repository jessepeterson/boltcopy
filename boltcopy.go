package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/boltdb/bolt"
)

type bucketFlags []string

func (b *bucketFlags) String() string {
	return strings.Join(*b, ",")
}

func (b *bucketFlags) Set(bucket string) error {
	*b = append(*b, bucket)
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"Usage:\n%s [flags] <intput.db> <output.db>\nFlags:\n",
			os.Args[0])
		flag.PrintDefaults()
	}

	var buckets bucketFlags
	var include = flag.Bool("i", false, "include provided buckets (i.e. white list behavior)")
	flag.Var(&buckets, "b", "name of bucket (can specify multiple)")

	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(3)
	}

	if _, err := os.Stat(flag.Arg(0)); os.IsNotExist(err) {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"file '%s' does not exist\n",
			flag.Arg(0))
		flag.Usage()
		os.Exit(3)
	}

	if _, err := os.Stat(flag.Arg(1)); err == nil {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"file '%s' exists, not overwriting\n",
			flag.Arg(1))
		flag.Usage()
		os.Exit(3)
	}

	if *include && len(buckets) < 1 {
		fmt.Fprintf(
			flag.CommandLine.Output(),
			"include-only mode (i.e. white list) but no buckets provided\n")
		flag.Usage()
		os.Exit(3)
	}

	idb, err := bolt.Open(flag.Arg(0), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer idb.Close()

	bucketList, err := genBucketCopyList(idb, buckets, *include)
	if err != nil {
		log.Fatal(err)
	}

	if len(bucketList) < 1 {
		fmt.Fprintln(os.Stderr, "error: no buckets to copy")
		os.Exit(4)
	}

	odb, err := bolt.Open(flag.Arg(1), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer odb.Close()

	odb.NoSync = true

	for _, b := range bucketList {
		err := copyBucket(idb, odb, b)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func genBucketCopyList(db *bolt.DB, buckets []string, includeOnly bool) ([]string, error) {
	bucketList := []string{}
	db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(nm []byte, _ *bolt.Bucket) error {
			found := false
			for _, b := range buckets {
				if b == string(nm) {
					found = true
					break
				}
			}
			if (includeOnly && found) || (!includeOnly && !found) {
				bucketList = append(bucketList, string(nm))
			}
			return nil
		})
		return nil
	})
	return bucketList, nil
}

func copyBucket(idb, odb *bolt.DB, bucket string) error {
	return idb.View(func(itx *bolt.Tx) error {
		ib := itx.Bucket([]byte(bucket))
		return odb.Update(func(otx *bolt.Tx) error {
			ob, err := otx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
			return ib.ForEach(func(k, v []byte) error {
				return ob.Put(k, v)
			})
		})
	})
}
