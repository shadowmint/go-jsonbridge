# JSON Bridge

A simple JSON bridge helper.

# Usage

    import "ntoolkit/jsonbridge"

    bridge = jsonbridge.New(istream, ostream)

    // Read from input stream
    bridge.Read()

    // Process all active messages
    for bridge.Len() > 0 {
        meta := Meta{}
        if bridge.As(&meta) == nil {
            if meta.type == 1 {
                msg := Message{}
                if bridge.As(&msg) == nil {
                    ...
                }
            } else if meta.type == 2 {
                msg := Message2{}
                if bridge.As(&msg) == nil {
                    ...
                }
            }
        }
    }
