- Feature Name: TAE Primary Key Index
- Status: draft
- Start Date: 2022-03-02
- Authors: [Zilong Zhou](https://github.com/zzl200012)
- Implementation PR:
- Issue for this RFC:

# Summary
Propose the first version of primary key index for TAE.

# Motivation
**TAE** (Transactional Analytical Engine) is a columnar storage engine supporting (1) snapshot isolation transaction and (2) efficient point query based on primary key. For (1), we need to verify if there are duplicate primary keys in the committing data; and for (2), we need to find the position of target row as fast as possible. Since that, a primary key index is needed, being both high-performance and easy-to-maintain.

# Technical Design
## High-level Design Decisions

### Granularity

As for the granularity of index, we divide the index into two categories, one is a table-level index, and the other is an index set composed of a series of partitioned indexes.

Many databases use a table-level B+Tree (or any other extensions like BwTree, ART, etc.) as primary key index. In theory that's pretty easy to use, and the query time is bounded as well. But in **TAE**, the data for one table consists of many segments, and each segment must be unordered first and then ordered. Compaction, merging, or splitting may take place afterwards, which makes the table-level index hard to maintain. Since that, the index of **TAE** should be more fine-grained, i.e. segment-level and block-level indexes, so that the lifetime of indexes would be bound to segments and blocks. Besides, if we choose a pure in-memory table-level index like many databases did (e.g. duckdb. we also considered fujimap as an option), the memory consumption would be large and unbounded. And for a disk-based table-level index, several I/Os per query are unavoidable, which results in poor performance. In a word, table-level index is a bad choice for our scenario.

<img src="https://github.com/zzl200012/docs-public/blob/main/seg-format.svg" height="60%" width="60%" />

The index is partitioned to each segments, and the layout of every segment is shown above (we only care about primary key index here, so other fields e.g. header, version nodes, checksum, etc. are all represented as "blk_*"). There are two types of segment in **TAE**, appendable or non-appendable. An appendable segment consists of at least one appendable block plus multiple non-appendable blocks. Appendable block index is an in-memory ART plus zonemap while the non-appendable one is a bloomfilter plus zonemap. For non-appendable segment, the index is a two-level structure, bloomfilter and zonemap respectively. As for bloomfilter, there are two options, a segment-based bloomfilter, or a block-based bloomfilter. The Segment-based is a better choice when the index can be fully resident in memory. The block-based is just like Rocksdb's approach.

As far as we know, Apache Kudu used similar approach as described above, to do both primary key deduplication and point query acceleration. In Kudu, every DiskRowSet (equivalent to our segment) has a bloom filter, plus a primary key index (using MassTree, an extension of B+Tree) which works the same as our zone map. Kudu manages all those structures with an LRU cache, and we do this as well.

### Maintenance

Currently there are two options for primary key index in **TAE**: in-memory or on-disk. For the zonemap, we can just let it resident in memory since it's really small. For ART, we know it's only available for appendable block, but only few of that exists simultaneously in **TAE**, the memory consumption is also limited. So the only concern is bloom filter, we should decide wether let it resident in memory or not.

<img src="https://github.com/zzl200012/docs-public/blob/main/lifetime.svg" height="60%" width="60%" />

In in-memory mode, the procedure is (A) -> (B) -> (C) -> (D) -> (F), while on-disk is (A) -> (B) -> (C) -> (D) -> (E). When a fresh appendable block generated, an ART and a zonemap would be attached with it. Every updates to this block would be applied to those two structures as well. When committing a block, the ART would be replaced with a bloom filter. Trick comes when segment is closed, i.e. the last block of the segment turns to non-appendable, and the segment turns to non-appendable as well. First, zonemaps separated into each block would be replaced with a segment zonemap, which would be more fine-grained and carry more informations. Then, for in-memory mode, bloom filters of each block would be replaced as well, with a segment bloom filter, otherwise just keep unchanged.

Notice that every structures' lifetime is bound to its host, block or segment, appendable or non-appendable. With the layout being changed, some of them are changed, and some are not. 

### Miscellaneous

#### Correctness

For primary key, there is only append or delete upon a segment. If delete occurs, we just ignore that, since false positive is allowed, only some more I/O cost on the underlying data is paid. When deletes are excessive, the whole segment/block would be re-constructed (via compaction), so we don't deal with it individually. For the changing part (i.e. appendable block), we assign an ART to get correct answer without any I/O.

#### Performance

Thanks to the fine-grained design, maintenance becomes easier, and memory consumption is controllable. So one last concern is performance, every time a request comes, we would iterate all segments until we **exactly** find the target or reach the end. "exactly" means once we get positive on a segment, we would dig into the underlying data and have a check whether it exists for real. Of course it would be super slow, so we have the following decisions:

1. Use better "bloom filter" optimized for static data, since our bloom filters are all built for immutable data, to gain better performance with lower false positive rate.
1. Use carefully designed "zone map", which means it's not a simple zonemap, but a more fine-grained secondary index like structure (e.g. imprints index), in order to avoid exact I/Os as more as possible.
1. In high-concurrency scenario, we can leverage such properties: (1) every segments are independent (2) inside a non-appendable segment/block, all fields are thread-safe since they are immutable, to build a query pipeline and boost our processing.

## Detailed Workflow

### Update

<img src="https://github.com/zzl200012/docs-public/blob/main/idx-update.svg" height="60%" width="60%" />

### Build

<img src="https://github.com/zzl200012/docs-public/blob/main/idx-build.svg" height="60%" width="60%" />

### Query

## Interface Definitions

# Future Works

