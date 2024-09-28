# Authorization-Enhanced Internet Mail Architecture

## Current Internet Mail Architecture (RFC5598)

![RFC5598 Sequence Diagram Part I.](../docs/resources/rfc5598_sd_I.puml)
<p class="figure">
    Fig.&nbsp;1.&emsp;Current Internet mail architecture (rfc5598), message flow part I.
</p>

![RFC5598 Sequence Diagram Part II.](../docs/resources/rfc5598_sd_II.puml)
<p class="figure">
    Fig.&nbsp;2.&emsp;Current Internet mail architecture (rfc5598), message flow part II.
</p>

## Authorization-Enhanced Internet Mail Architecture

![Authorization-Enhanced RFC5598 Sequence Diagram Part I.](../docs/resources/authorization-enhanced_rfc5598_I_sd.puml)
<p class="figure">
    Fig.&nbsp;3.&emsp;Authorization-Enhanced Internet mail architecture, message flow part I.
</p>

![Authorization-Enhanced RFC5598 Sequence Diagram Part II.](../docs/resources/authorization-enhanced_rfc5598_II_sd.puml)
<p class="figure">
    Fig.&nbsp;4.&emsp;Authorization-Enhanced Internet mail architecture, message flow part II.
</p>

## Placeholder Message

An example of placeholder message in JSON format with external bodies accessible via content-addressed URIs.

```json
{
  "headers":
    {
      "X-Author-RS-URL": "https://foo.com/rs",
      "X-Recipient-RS-URL": "https://bar.com/rs",
      "From": "Alice Sanders <alice@foo.com>",
      "Subject": "Meeting",
      "To": "Bob Sanders <bob@bar.com>",
      "Cc": "Carol <carol@bar.com>, Daniel <dan@bar.com>",
      "Date": "Tue Sep 19 20:52:05 CEST 2023",
      "Message-ID": "<b07d0cdf-c6f4-4f67-b24c-cc847a4c2df4@foo.com>",
      "X-Thread-ID": "<68fb9177-6853-466a-8f7d-c96fbb885f81@foo.com>",
      "Content-Type": "multipart/mixed"
    },
  "parts":
    [
      {
        "headers": { "Content-Type": "multipart/alternative" },
        "parts":
          [
            {
              "headers":
                {
                  "Content-Disposition": "inline",
                  "Content-ID": "<aSQnmlBT6RndpDnwTSStJUVhlh9XL9_y2QXX42NhKuI>",
                  "Content-Type":
                    [
                      "message/external-body; access-type='x-content-addressed-uri'; hash-algorithm='sha256'; size='42'",
                      "text/plain; charset=UTF-8"
                    ]
                }
            },
            {
              "headers":
                {
                  "Content-Disposition": "inline",
                  "Content-ID": "<Y_ION3g8WQuqGzhsDlVrhAgQ0D7AbXu9T-HSv3w--zY>",
                  "Content-Type":
                    [
                      "message/external-body; access-type='x-content-addressed-uri'; hash-algorithm='sha256'; size='109'",
                      "text/html; charset=UTF-8"
                    ]
                }
            }
          ]
      },
      {
        "headers": { "Content-Type": "multipart/mixed" },
        "parts":
          [
            {
              "headers":
                {
                  "Content-Disposition": "attachment; filename='logo.svg'",
                  "Content-ID": "<1pzyqfFWbfhJ3hrydjL9jO9Qgeg70TgZQ_zpOkt4HOU>",
                  "Content-Type":
                    [
                      "message/external-body; access-type='x-content-addressed-uri'; hash-algorithm='sha256'; size='52247'",
                      "image/svg+xml"
                    ]
                }
            },
            {
              "headers":
                {
                  "Content-Disposition": "attachment; filename='Minutes.pdf'",
                  "Content-ID": "<6G6Mkapa3-Om7B6BVhPUBEsCLP6t6LAVP4bHxhQF5nc>",
                  "Content-Type":
                    [
                      "message/external-body; access-type='x-content-addressed-uri'; hash-algorithm='sha256'; size='153403'",
                      "application/pdf"
                    ]
                }
            }
          ]
      }
    ]
}
```

