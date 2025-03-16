# pbr [![CI](https://github.com/pchchv/pbr/workflows/CI/badge.svg)](https://github.com/pchchv/pbr/actions?query=workflow%3ACI+event%3Apush) [![Go Report Card](http://goreportcard.com/badge/github.com/pchchv/pbr)](https://goreportcard.com/report/github.com/pchchv/pbr) [![Godoc Reference](https://pkg.go.dev/badge/github.com/pchchv/pbr)](https://pkg.go.dev/github.com/pchchv/pbr)

Package `pbr` is a low-level reader for [protocol buffers](https://developers.google.com/protocol-buffers) encoded data in Golang.  
Main feature is the support for lazy/conditional decoding of fields.

This package can help decoding performance in two ways:

1. fields can be conditionally decoded, skipping fields that are not needed for a specific use case

2. decoding directly into specific types or performing other transformations, where additional state can be skipped by manually decoding into types directly.

**Note:** using [gogoprotobuf](https://github.com/gogo/protobuf) is still faster.   
**Note:** Writing code with this package is like writing an auto-generated protobuf decoder and is very time consuming. It should be used only in specific cases and for stable protobuf definitions.

## Usage

First, the encoded protobuf data is used to initialize a new Message. Then the fields are searched by reading or skipping them.

```go
msg := pbr.New(encodedData)
for msg.Next() {
    switch msg.FieldNumber() {
    case 1: // an int64 type
        v, err := msg.Int64()
        if err != nil {
            // handle
        }
    case 3: // repeated number types can be returned as a slice
        ids, err := msg.RepeatedInt64(nil)
        if err != nil {
            // handle
        }
    case 2: // for more control repeated+packed fields can be read using an iterator
        iter, err := msg.Iterator(nil)
        if err != nil {
            // handle
        }

        userIDs := make([]UserID, 0, iter.Count(pbr.WireTypeVarint))
        for iter.HasNext() {
            v, err := iter.Int64()
            if err != nil {
                // handle
            }

            userIDs = append(userIDs, UserID(v))
        }
    default:
        msg.Skip() // required if value not needed.
    }
}

if msg.Error() != nil {
    // handle
}
```

After calling `Next()` you **must** call an accessor function (`Int64()`, `RepeatedInt64()`, `Iterator()`, etc.) or `Skip()` to ignore the field. All these functions, including `Next()` and `Skip()`, must not be called twice in a row.

### Value Accessor Functions

There is an accessor for each one the protobuf [scalar value types](https://developers.google.com/protocol-buffers/docs/proto#scalar).

There is a corresponding set of functions for repeated fields, e.g. `RepeatedInt64(buf []int64) ([]int64, error)`. Repeated fields may or may not be packed, so a predefined buffer variable should be passed when called. E.g.:

```go
var ids []int64
msg := pbr.New(encodedData)
for msg.Next() {
    switch msg.FieldNumber() {
    case 1: // repeated int64 field
        var err error
        ids, err = msg.RepeatedInt64(ids)
        if err != nil {
            // handle
        }
    default:
        msg.Skip()
    }
}

if msg.Error() != nil {
    // handle
}
```

If the ids are 'packed', `RepeatedInt64()` will be called once. If the ids are simply repeated `RepeatedInt64()` will be called N times, but the resulting array of ids will be the same.

For more control over the values in a packed, repeated field use an Iterator.  
See above for an example.

### Decoding Embedded Messages

Embedded messages can be handled recursively, or the raw data can be returned and decoded using a standard/auto-generated `proto.Unmarshal` function.

```go
msg := pbr.New(encodedData)
for msg.Next() {
    fn := msg.FieldNumber()
    // use pbr recursively
    if fn == 1 && needFieldNumber1 {
        embeddedMsg, err := msg.Message()
        for embeddedMsg.Next() {
            switch embeddedMsg.FieldNumber() {
            case 1:
                // do something
            default:
                embeddedMsg.Skip()
            }
        }
    }

    // if you need the whole message decode the message in the standard way.
    if fn == 2 && needFieldNumber2 {
        data, err := msg.MessageData()
        v := &ProtoBufThing()
        err = proto.Unmarshal(data, v)
    }
}
```

## Larger Example
Start with a customer message with embedded orders and items, need to count only the number of items in open orders.

```protobuf
message Customer {
  required int64 id = 1;
  optional string username = 2;
  repeated Order orders = 3;
  repeated int64 favorite_ids = 4 [packed=true];
}

message Order {
  required int64 id = 1;
  required bool open = 2;
  repeated Item items = 3;
}

message Item {
  // a big object
}
```

Sample Code:

```go
var openCount, itemCount, favoritesCount int
customer := pbr.New(data)
for customer.Next() {
    switch customer.FieldNumber() {
    case 1: // id
        id, err := customer.Int64()
        if err != nil {
            panic(err)
        }
        _ = id // do something or skip this case if not needed
    case 2: // username
        username, err := customer.String()
        if err != nil {
            panic(err)
        }
        _ = username // do something or skip this case if not needed
    case 3: // orders
        var open bool
        var count int
        orderData, _ := customer.MessageData()
        order := pbr.New(orderData)
        for order.Next() {
            switch order.FieldNumber() {
            case 2: // open
                v, _ := order.Bool()
                open = v
            case 3: // item
                count++

                // we're not reading the data but we still need to skip it.
                order.Skip()
            default:
                // required to move past unneeded fields
                order.Skip()
            }
        }

        if open {
            openCount++
            itemCount += count
        }
    case 4: // favorite ids
        iter, err := customer.Iterator(nil)
        if err != nil {
        	panic(err)
        }

        // typically this section would only be run once but it is valid
        // protobuf to contain multiple sections of repeated fields that should be concatenated together
        favoritesCount += iter.Count(pbr.WireTypeVarint)
    default:
        // unread fields must be skipped
        customer.Skip()
    }
}

fmt.Printf("Open Orders: %d\n", openCount)
fmt.Printf("Items:       %d\n", itemCount)
fmt.Printf("Favorites:   %d\n", favoritesCount)

// Output:
// Open Orders: 2
// Items:       4
// Favorites:   8
```

## Wire Type Start Group and End Group

Groups are an old type of protobuf wires that have been deprecated for a long time.
They work like parentheses, but do not contain any information about the length of the data.
Therefore, their contents cannot be effectively skipped. Only the start and end group indicators can be read and skipped like any other field. This will cause the data to be read without parentheses, whatever that means in practice. To get the raw protobuf data within a group, try the following:

```go
var groupFieldNum = 123
var groupData []byte
msg := New(data)
for msg.Next() {
    if msg.FieldNumber() == groupFieldNum && msg.WireType() == WireTypeStartGroup {
        start, end := msg.Index, msg.Index
        for msg.Next() {
            msg.Skip()
            if msg.FieldNumber() == groupFieldNum && msg.WireType() == WireTypeEndGroup {
                break
            }
            end = msg.Index
        }
        // groupData would be the raw protobuf encoded bytes of the fields in the group.
        groupData = msg.Data[start:end]
    }
}
```