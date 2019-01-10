Package rangearray implements indexing and searching in a semi-dense
array.  In particular, a semi-dense array has runs of consecutive
entries with potentially large (and irregular) separation between
those runs.

Note that the data structure is optimized for insertions only at the
end; insertions before the end can be very expensive.  A search over
n runs takes O(log(n)) time.
