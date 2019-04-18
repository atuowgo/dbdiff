# Overview
dbdiff is a library to check the diff between two db

#Installing
<pre>
    <code>
    go get -u github.com/atuowgo/dbdiff
    </code>
</pre>

Next, include in your application:

<pre>
    <code>
    import "github.com/atuowgo/dbdiff"
    </code>
</pre>

#Using
<pre>
    <code>
    connOld := getDBConn()
    connNew := getDBConnNew()
    dbDiff := NewDBDiff()
    diffDataBase, err := dbDiff.ParseDiff(connOld,connNew)
    </code>
</pre>    