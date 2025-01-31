# System Instructions

You are an expert in OAuth 2.0 and email system. You are an experienced frontend and backend developer. You know how to code in these languages (a) Golang, SQL (b) HTML, Javascript/Typescript, CSS and you know how to use React framework with styled-components. You can write API documentation using OpenAPI 3.1 specification. You know how to design a synchronization protocol to synchronize data between backend services and a frontend application.

# The Federizer Project

We aim to create an Internet mail system, similar to the current email system, that facilitates the transmission of various types of data, such as messages, documents, books, photos, podcasts, and videos, using a set of OAuth 2.0 mechanisms.

## Logic

### Placeholder Message and Confined External Resources

Each email entity consists of a Placeholder Message and its associated External Resources—the message bodies stored within the mailbox.

The External Resources are confined to the individuals listed in the "From," "To," "Cc," and "Bcc" headers of the Placeholder Message. This implies that the Placeholder Message functions as an access control list, granting these individuals access to the External Resources.

Owners of the External Resources can send each recipient a copy of the Placeholder Message, signed by their agent. This action notifies the recipients' agents to fetch the corresponding External Resources.

### Contextual Discharge

To discharge the External Resources to additional recipients specified in the "Forwarded-From" and "Forwarded-To" headers of the Placeholder Message, a chain of signed Placeholder Message Envelope Headers must be presented. As the Placeholder Message passes through each mailbox service, it accumulates a contextual sequence of signed Placeholder Message Envelope Headers. This chain enables the agents of these recipients to access and fetch the External Resources.

Upon fetching the External Resources, each recipient gains ownership of the copies stored in their mailbox.

## Concept

### Acronyms

For the sake of brevity of this document, the following list of acronyms will be used:

* AS: Authorization Server
* RS: Resource Server
* MTA: Mail Transfer Agent
* MBX: Mailbox
* APP: Webmail Application

### Components

Suppose we have two trust domains, example.com and example.net, and two users: (a) Alice with the email address alice@example.com, and (b) Bob with the email address bob@example.net. These domains and users do not have any pre-established trust relationship, nor do we intend to create any permanent trust relationship (e.g. through federation) between the trust domains.

We have the following components:

* Authorization Server (AS1): Operates under the example.com trust domain.

* Authorization Server (AS2): Operates under the example.net trust domain.

* Client (Client1): Operates under the example.com trust domain and serves as Alice's Webmail application (APP1).

* Client (Client2): Operates under the example.com trust domain and serves as Alice's Mail Transfer Agent (MTA1)

* Client (Client3): Operates under the example.net trust domain and serves as Bob's Mail Transfer Agent (MTA2)

* Client (Client4): Operates under the example.net trust domain and serves as Bob's Webmail application (APP2).

* Resource Server (RS1): Operates under the example.com trust domain and serves as Alice's Mailbox (MBX1).

* Resource Server (RS2): Operates under the example.net trust domain and serves as Bob's Mailbox (MBX2).

* MBX1 and MTA1: We will refer to them collectively as the MBX1/MTA1 entity, that acts in a dual role and functions as both a service and an agent:

  * Acts as an RS1 with respect to Client1 or Client3.

  * Acts as a Client2 with respect to RS2.

* MBX2 and MTA2: We will refer to them collectively as the MBX2/MTA2 entity, that acts in a dual role and functions as both a service and an agent:

  * Acts as an RS2 with respect to Client4 or Client2.

  * Acts as a Client3 with respect to RS1.

The MTA allows bidirectional data transfer, it can send and receive resources using the POST and GET http methods. That allows to deliver Alice's resources to Bob in two different ways:

1. Alice's MTA1 sends her resources to Bob's MBX2.

2. Alice temporarily shares her resources to Bob on her MBX1 and sends a notification message (a Placeholder Message, see Fig. 1) via MTA1 to Bob's MBX2/MTA2 service. Next, the Bob's MTA2 fetches the resources from Alice's MBX1/MTA1 service.

Taking into account both ways of transferring Alice's resources to Bob and highlight the symmetry of the system (Alice can deliver her resources to Bob in two ways; Bob can deliver his resources to Alice in two ways) we can summarize the concept of dual roles as follows:

* the MBX1/MTA1 acts as RS1 when Alice posts/gets her resources via APP1.

* the MBX1/MTA1 acts as Client2 when transferring Alice's resources to Bob's MBX2.

* the MBX1/MTA1 acts as RS1 when Bob's MTA2 fetches Alice's shared resources from Alice's MBX1.

* the MBX2/MTA2 acts as RS2 when Bob posts/gets his resources via APP2.

* the MBX2/MTA2 acts as Client3 when transferring Bob's resources to Alice's MBX1.

* the MBX2/MTA2 acts as RS2 when Alice's MTA1 fetches Bob's shared resources from Bob's MBX2.

### OAuth 2.0 Flow

We utilize two OAuth 2.0 grant types:

1. Authorization Code Grant with PKCE: Facilitates user-driven interactions (e.g., Alice or Bob accessing their respective mailboxes using Webmail applications).

2. Token Exchange (RFC 8693) for Email Transfer with JWT Assertion and the Demonstrating Proof of Possession (DPoP) mechanism using the "eh" claim: Authorizes automated service-to-service interactions (e.g., Alice's MTA1 accessing Bob's Mailbox MBX2). We call this complex authorization mechanism Cross-Domain Authorization Grant and will describe it in more detail later.

### Client Registration

We use the following OAuth 2.0 Client Registration Schema:

* The APP1 is registered at the AS1 as a public client.

* The APP2 is registered at the AS2 as a public client.

* The MBX1/MTA1 is registered at the AS1 as a confidential client.

* The MBX2/MTA2 is registered at the AS2 as a confidential client.

### Service Discovery

In the OAuth 2.0-based internet mail architecture, the mailbox (locator) and the email address (user identity identifier) are separated. This separation allows users from the example.com trust domain to have their mailboxes hosted in another domain, such as example.edu. This architecture always requires two DNS records: one for the identity provider and another for the storage provider. However, both providers may operate under the same trust domain, as we will assume in the following description.

We implement service discovery using these DNS records:

1. `_federizer._as._tcp.example.com. 3600 IN SRV 10 5 443 as1.example.com.`

   * Points either: (a) directly to the authorization service, or (b) to a redirect service that redirects http requests to the URL of the authorization service (e.g., example.com/as1), see Note 1. This is used by the MBX2/MTA2 service to verify that Alice's AS1 is the corresponding JWT Assertion issuer.

2. `_federizer._rs._tcp.example.com. 3600 IN SRV 10 5 443 rs1.example.com.`

   * Points either: (a) directly to the resource server, or (b) to a redirect service that redirects http requests to the URL of the resource server (e.g., example.com/rs1), see Note 1. This is used by the MBX2/MTA2 agent to discover the MBX1/MTA1 service and by the MBX2/MTA2 service to verify that the MBX1/MTA1 agent is the correct client making the request, see Note 2.

3. `_federizer._as._tcp.example.net. 3600 IN SRV 10 5 443 as2.example.net.`

   * Points either: (a) directly to the authorization service, or (b) to a redirect service that redirects http requests to the URL of the authorization service (e.g., example.net/as2), see Note 1. This is used by the MBX1/MTA1 service to verify that Bob's AS2 is the corresponding JWT Assertion issuer.

4. `_federizer._rs._tcp.example.net. 3600 IN SRV 10 5 443 rs2.example.net.`

   * Points either: (a) directly to the resource server, or (b) to a redirect service that redirects http requests to the URL of the resource server (e.g., example.net/rs2), see Note 1. This is used by the MBX1/MTA1 agent to discover the MBX2/MTA2 service and by the MBX1/MTA1 service to verify that the MBX2/MTA2 agent is the correct client making the request, see Note 2.

Note 1: DNS SRV records cannot point to a URL with a path, while using a URL with a specific path is common for Authorization Servers or Resource Servers. A trusted redirect service can run on the hostname specified in the SRV record and redirect (e.g., a 302 Found) http requests to the actual URL of the Authorization Server or Resource Server, including the path.

Note 2: We utilize the dual role feature of MBX1/MTA1 and MBX2/MTA2 entities to provide metadata about Client2 and Client3. When each MTA is registered with its respective AS as a confidential client, the sibling MTX resource server can supply specific metadata about its MTA, including the MTA client_id. This client_id, also included in the JWT Assertion as an "azp" claim, is used to verify that the correct client is making the request. We refer to this mechanism as a DNS-bound client identity via the sibling SRV record. This allows for the use of DNS records to specify an OAuth 2.0-based internet mail hosting provider.

### Conceptual Example

* The External Resources owned by the author stored on the RS of the origin mailbox service are temporarily shared with recipients by creating a Placeholder Message which also acts as an access control list. Following a successful sharing process, using the Cross-Domain Authorization Grant, the Placeholder Message is sent to each recipient through the MTA operating within the origin trust domain. This Placeholder Message stores SHA-256 digests of the referenced External Resources in its `Content-ID` headers (see Figure 1 for an example of a Placeholder Message).

* The MTA operating within the destination trust domain using the Cross-Domain Authorization Grant attempts to fetch the External Resources from the RS of the origin mailbox service. After successful authorization, the External Resource is fetched using a digest value of the Content-ID header from the Placeholder Message body as an identifier and is stored on the RS of the destination mailbox service. That means that each recipient becomes the owner of the corresponding copy of the referenced External Resource (a Signed Placeholder Message can prove it, as explained later—TBD), which they can download and use, or send to other recipients. Finally, the Webmail application downloads the relevant data from the RS of the destination mailbox service and reconstructs the original message according to the Placeholder Message source.

#### Draft Placeholder Message

```yaml
# Envelope
- headers:
    - From: Alice <alice@example.com>
    - Subject: Vacation photo
    - To: Bob <bob@example.net>
    - Message-ID: <b07d0cdf-c6f4-4f67-b24c-cc847a4c2df4@example.com>
    - X-Thread-ID: <68fb9177-6853-466a-8f7d-c96fbb885f81@example.com>
    - Content-Type: multipart/mixed
# Body
- parts:
    - headers:
        - Content-Type: multipart/alternative
      parts:
        - headers:
            - Content-Disposition: inline
            - Content-ID: <aSQnmlBT6RndpDnwTSStJUVhlh9XL9_y2QXX42NhKuI>
            - Content-Type:
              - message/external-body; access-type='x-content-addressed-uri';
                hash-algorithm='sha256'; size='42'
              - text/plain; charset=UTF-8
        - headers:
            - Content-Disposition: inline
            - Content-ID: <Y_ION3g8WQuqGzhsDlVrhAgQ0D7AbXu9T-HSv3w--zY>
            - Content-Type:
              - message/external-body; access-type='x-content-addressed-uri';
                hash-algorithm='sha256'; size='109'
              - text/html; charset=UTF-8
    - headers:
        - Content-Type: multipart/mixed
      parts:
        - headers:
            - Content-Disposition: attachment; filename='Venice.png'
            - Content-ID: <1pzyqfFWbfhJ3hrydjL9jO9Qgeg70TgZQ_zpOkt4HOU>
            - Content-Type:
              - message/external-body; access-type='x-content-addressed-uri';
                hash-algorithm='sha256'; size='3520247'
              - image/png
```
Fig. 1. An example of a draft Placeholder Message in YAML format.

In the example above, all headers are part of the original email composition. Additional headers, such as `Date`, `Received`, `Forwarded-From`, and `Forwarded-To`, are added sequentially as the Placeholder Message passes through each mailbox service. The most recent header appears at the top of the list (see Figure 2).

#### Signed Placeholder Message Envelope Headers

```yaml
# Envelope
- headers:
    - Message-ID: <b07d0cdf-c6f4-4f67-b24c-cc847a4c2df4@example.com>
    - JWT-Assertion-Digest: <PpN6Qr16xzNhDx_stziNtYEUXHlJ8v1OVxGELWr0tuY> # last hop JWT Assertion Digest
    - Received: from example.com by example.net; Sun Dec 22 20:58:14 CEST 2024
- jwt-assertion:
    header:
      alg: RS256
      typ: JWT
    payload:
      iss: https://as1.example.com/
      azp: client2.example.com
      scope: internet_mail
      cnf:
        jkt: 0ZcOCORZNYy-DWpqq30jZyJGHTN0d2HglBV3uiguA4I # DPoP JWK Thumbprint (here superfluous)
      eh:                                                # last hop Envelope Headers
        eht: tKKw0-Munn-PQKcmE_nuF32eziwCS0FBtUz_fJMva5E # Envelope Headers Thumbprint
        ehl:                                             # Named list of hashed Envelope Headers
          - Message-ID
          - From
          - Recipients-Digest
          - Body-Digest
    signature: GHSqC5F1H6D7SkN0JwhVw6aKt8TNuBWKoKXv19HlkMJP6dQPfw9u6LgxEyRBjt08STyqIAqSPuu7tzWzL_efxYi1j_9ABteTcHDC9lJRNAwp9TNmg9K4uXtdSw57K3vdzDfeCxcfwhJS_t3wz1vtyF3pvsdkJ5G81qxZ8Gh8n5hZsadzknL6yCeL_x9K_ykDFJl_TVBUf7Q1ImyEUaCx04DA-eZv85xlWkG0NeutVQ6epeB-XAuFxImaTzM-kP66YXDNeK1ohtoS99Qk9S_e0cQZV7dl3wtxSCf3A-sqTrh0nDPn8D71annnXBsqUWJsD4QERed-lKnGagbzJn57mw
- headers:
    - Date: Sun Dec 22 20:49:35 CEST 2024
    - Recipients-Digest: <tSutyHKcSxIt_2VxOi1mYtU5h9ffs-AHK9d6ffwxmle> # To, Cc, Bcc headers digest
    - Body-Digest: <nZURGvgk4xoy6-aI6dna5ddskq5ud_GyI7u96hkxYh4>
    - From: Alice <alice@example.com>
    - Subject: Vacation photo
    - To: Bob <bob@example.net>
    - Message-ID: <b07d0cdf-c6f4-4f67-b24c-cc847a4c2df4@example.com>
    - X-Thread-ID: <68fb9177-6853-466a-8f7d-c96fbb885f81@example.com>
    - Content-Type: multipart/mixed
```
Fig. 2. An example of a signed Placeholder Placeholder Message Envelope Headers in YAML format at the end of the send/receive process.

The `Date` header means that the message has been sent and is immutable—without the `Date` header it is a draft message.
The `Recipients-Digest` is a digest of the concatenated values of the `To`, `Cc`, `Bcc` headers, or of the single value of the`Forwarded-To` header. It hides the real recipients from being exposed in the authorization server.
The `Body-Digest` header binds the message body to the envelope headers. Body binding through a body digest involves generating a SHA-256 hash of the message body.
The `eht` claim binds the envelope headers to the JWT Assertion.

#### Chain of Signed Placeholder Message Envelope Headers

```yaml
# Envelope
- headers:
    - Message-ID: <b07d0cdf-c6f4-4f67-b24c-cc847a4c2df4@example.com>
    - JWT-Assertion-Digest: <tVaz0gjAdWeLn5JIz8g8NUNVnJL2ZAPnEcapTNM0zc7> # last hop JWT Assertion Digest
    - Received: from example.net by example.org; Sun Dec 22 21:01:34 CEST 2024
- jwt-assertion:
    header:
      alg: RS256
      typ: JWT
    payload:
      iss: https://as2.example.net/
      azp: client3.example.net
      scope: internet_mail
      cnf:
        jkt: gFSUE65ghV8FBld-0J85pwB4TOKKqEfMvvd8Rm8Hy9G # DPoP JWK Thumbprint (here superfluous)
      eh:                                                # last hop Envelope Headers
        eht: zTvzWxYg9zuqZHvZhg6RvdpqPpLvtXMS-l8QUehswJP # Envelope Headers Thumbprint
        ehl:                                             # Named list of hashed Envelope Headers
          - Message-ID
          - Forwarded-From
          - Recipients-Digest
          - JWT-Assertion-Digest
    signature: Drs2jP8vasca5qCDmmB-spUuPqrdAq1oGVaObRC_kTgdlx48hpDJcYZ18mxBTGJ_cnoRQ1gBdlHaKnnM0ir6QPEt3qDpuU7mT9v6Gj84GhsXfbsGGpIWPL4ZNPI-u8TNeZWMZc0VkVNx1FswIOMQNLLQsPZ1gPPGLnT7LTWzyib8leE779itPLyI6jMrhrYxOfbGRIb4V9jIiRxsbTKUH1tNdZ9Hy5PC75u24-YiUNup8tF4O9JStzoe3K6QJG0jHjhVRJ_efqt6tBPdl2XfjZOEr_1cMIFsk6iV_2GiqqacHRKzC0JoBp8Djxyp_DaD8OP8m3AbyHP5Ial9s7zhhg
- headers:
    - Recipients-Digest: <i0LnZ-Ju2O3_1E2qNFRg2T7DwtMePhvQnWqTQO9ZpWK> # Forwarded-To header digest    
    - Forwarded-To: Bobby <bobby@example.org>
    - Forwarded-From: Bob <bob@example.net>
    - Message-ID: <b07d0cdf-c6f4-4f67-b24c-cc847a4c2df4@example.com>
    - JWT-Assertion-Digest: <PpN6Qr16xzNhDx_stziNtYEUXHlJ8v1OVxGELWr0tuY> # last hop JWT Assertion Digest
    - Received: from example.com by example.net; Sun Dec 22 20:58:14 CEST 2024
- jwt-assertion:
    header:
      alg: RS256
      typ: JWT
    payload:
      iss: https://as1.example.com/
      azp: client2.example.com
      scope: internet_mail
      cnf:
        jkt: 0ZcOCORZNYy-DWpqq30jZyJGHTN0d2HglBV3uiguA4I # JWK Thumbprint (here superfluous)
      eh:                                                # last hop Envelope Headers
        eht: tKKw0-Munn-PQKcmE_nuF32eziwCS0FBtUz_fJMva5E # Envelope Headers Thumbprint
        ehl:                                             # Named list of hashed Envelope Headers
          - Message-ID
          - From
          - Recipients-Digest
          - Body-Digest
    signature: GHSqC5F1H6D7SkN0JwhVw6aKt8TNuBWKoKXv19HlkMJP6dQPfw9u6LgxEyRBjt08STyqIAqSPuu7tzWzL_efxYi1j_9ABteTcHDC9lJRNAwp9TNmg9K4uXtdSw57K3vdzDfeCxcfwhJS_t3wz1vtyF3pvsdkJ5G81qxZ8Gh8n5hZsadzknL6yCeL_x9K_ykDFJl_TVBUf7Q1ImyEUaCx04DA-eZv85xlWkG0NeutVQ6epeB-XAuFxImaTzM-kP66YXDNeK1ohtoS99Qk9S_e0cQZV7dl3wtxSCf3A-sqTrh0nDPn8D71annnXBsqUWJsD4QERed-lKnGagbzJn57mw
- headers:
    - Date: Sun Dec 22 20:49:35 CEST 2024
    - Recipients-Digest: <tSutyHKcSxIt_2VxOi1mYtU5h9ffs-AHK9d6ffwxmle> # To, Cc, Bcc headers digest
    - Body-Digest: <nZURGvgk4xoy6-aI6dna5ddskq5ud_GyI7u96hkxYh4> # Body Digest
    - From: Alice <alice@example.com>
    - Subject: Vacation photo
    - To: Bob <bob@example.net>
    - Message-ID: <b07d0cdf-c6f4-4f67-b24c-cc847a4c2df4@example.com>
    - X-Thread-ID: <68fb9177-6853-466a-8f7d-c96fbb885f81@example.com>
    - Content-Type: multipart/mixed
```
Fig. 3. An example of chained, signed Placeholder Message Envelope Headers in YAML format at the end of the send/receive/forward/receive process.

## Digital Assets

The user can digitally sign the Placeholder Message, referred to as a PGP Signed Placeholder Message, either as a whole (the email body) or in parts by group of headers (`Content-Disposition`, `Content-ID`, `Content-Type`) of the individual External Resources using the `application/pgp-signature` protocol. The signature transforms the External Resource into a digital asset that has verifiable authenticity, integrity, and data origin, provided the resource includes provenance metadata—TBD.

## Resource(s) Delivery Scenario

We describe the delivery of Alice's resources to Bob.

Alice wants to send a vacation photo to Bob. She creates a new Placeholder Message in YAML format using the APP1 compose form. Alice opens the compose form and fills in the "To", "Subject", and (optionally) "Text" fields. While composing the message, she uploads the photo to the MBX1/MTA1 service to add it as an External Resource to the Placeholder Message. After a successful upload, the service returns a digest of the uploaded photo. This digest is then added to the Placeholder Message as the value of the `Content-ID` header and an External Resource with the `Content-Disposition: attachment` header. The content of the "Text" field is posted as a resource to the MBX1/MTA1 service that returns a digest of the posted text. This digest is then added to the Placeholder Message as the value of the `Content-ID` header and an External Resource with the `Content-Disposition: inline` header. Finally, to initiate the delivery process, Alice presses the "Send" button, posting the Placeholder Message to the MBX1/MTA1 service.

The MBX1/MTA1 agent sends the Placeholder Message to the MBX2/MTA2 service, using the Cross-Domain Authorization Grant. After completing the validation process, the MBX2/MTA2 agent fetches the External Resources from the MBX1/MTA1 service using the Cross-Domain Authorization Grant and notifies the APP2 Webmail application that it has received a new email.

Upon receiving the notification, Bob's APP2 downloads the Placeholder Message from the MBX2/MTA2 service. Using the information contained within, APP2 downloads the corresponding External Resources with the `Content-Disposition: inline` header (the "Text" field) to accurately reconstruct and display the content of the original Placeholder Message (including the content of the "Text" field). After this  delivery process, Bob is able to download the External Resource with the `Content-Disposition: attachment` header (the vacation photo) identified by its digest from the MBX2/MTA service.

<div style="break-after:page"></div>

## Sequence Diagrams

### OAuth 2.0 Authorization Code Grant with PKCE

This diagram illustrates the sequence of interactions in an OAuth 2.0 Authorization Code Grant flow enhanced with Proof Key for Code Exchange (PKCE). It is used in the interaction of Alice's APP1 with Alice's MBX1/MTA1 service and also in the interaction of Bob's App2 with Bob's MBX2/MTA2 service.

#### To authorize a public client to interact with the backend service:

* We will use OAuth 2.0 Authorization Code Grant with PKCE (RFC 7636).
* The client should be registered at the authorization server as a public client.

##### Sequence Diagram

```plantuml
@startuml
title OAuth 2.0 Authorization Code Grant with PKCE

actor User as "User/Browser"
participant Client as "Client Application"
participant "Authorization Server" as AuthServer
participant "Resource Server" as ResourceServer

User -> Client: Initiate authorization request
activate Client
create "Code Verifier" as CodeVerifier
Client -> CodeVerifier: Generate code_verifier
activate CodeVerifier
CodeVerifier --> Client: code_verifier
deactivate CodeVerifier
create "Code Challenge" as CodeChallenge
Client -> CodeChallenge: Generate code_challenge from code_verifier (SHA256)
activate CodeChallenge
CodeChallenge --> Client: code_challenge
deactivate CodeChallenge
Client -> AuthServer: Authorization Request\n(response_type=code, client_id, redirect_uri,\nscope, state, code_challenge, code_challenge_method=S256)
activate AuthServer
AuthServer -> User: Authentication Request
activate User
User -> AuthServer: Authenticate
AuthServer -> User: Consent Request (if required)
User -> AuthServer: Grant Consent (if required)
AuthServer --> Client: Authorization Code\n(code, state)
deactivate User
Client -> AuthServer: Token Request\n(grant_type=authorization_code, code, redirect_uri,\nclient_id, code_verifier)
AuthServer -> CodeVerifier: Retrieve code_verifier
activate CodeVerifier
CodeVerifier --> AuthServer: code_verifier
deactivate CodeVerifier
AuthServer -> CodeChallenge: Generate code_challenge from code_verifier (SHA256)
activate CodeChallenge
CodeChallenge --> AuthServer: code_challenge
deactivate CodeChallenge
AuthServer -> AuthServer: Verify code_challenge
AuthServer --> Client: Access Token, Refresh Token (optional),\nID Token (if OpenID Connect)\n(access_token, token_type, expires_in, refresh_token, id_token)
deactivate AuthServer
Client -> ResourceServer: Access Protected Resource\n(access_token)
activate ResourceServer
ResourceServer --> Client: Protected Resource
deactivate ResourceServer
deactivate Client

@enduml
```

The sequence diagram of this standard flow is self-explanatory.

<div style="break-after:page"></div>

### OAuth 2.0 Cross-Domain Authorization Grant

We need to replace the DomainKeys Identified Mail (DKIM) email authentication method by the Cross-Domain Authorization Grant, which we present as a JWT Assertion to Mailbox services running in different trust domains.

#### To authorize the transfer of the Placeholder Message:

* We will use the OAuth 2.0 Token Exchange grant type (RFC 8693).
* We will exchange the Access Token for a JWT Assertion.
* We will use a JWT Assertion with the "eht" claim (a SHA-256 thumbprint of the set of Placeholder Message Envelope Headers that serves as a digital signature). This assertion embodies the DKIM signature.
* We will use the DPoP mechanism to bind the JWT Assertion to the client.
* The MBX1/MTA1 agent should be registered at the AS1 authorization server as a confidential client.

##### Sequence Diagram

```plantuml
@startuml
title OAuth 2.0 Token Exchange for transferring Placeholder Messages using JWT Assertion and DPoP
participant "Client1" as Client1
participant "MBX1/MTA1" as Entity1
participant "AS1 Authorization Server" as AuthServer
participant "MBX2/MTA2" as Entity2

activate Client1
Client1 -> Entity1: Send Placeholder Message\n (includes Access Token)
activate Entity1

Entity1 -> Entity1: Generate DPoP Proof (using new key pair)
Entity1 -> Entity1: Create Body-Digest=<SHA256(Placeholder Message Body)> header\n and add it to the Placeholder Message
Entity1 -> Entity1: Create Recipients-Digest=<SHA256(To||Cc||Bcc)> header\n and add it to the Placeholder Message
Entity1 -> Entity1: Create Date=<Date> header and add it to the Placeholder Message
Entity1 -> Entity1: Store the Placeholder Message in the store, using the Message-ID header value as\n an identifier, <Stored Placeholder Message>=<Placeholder Message>
Entity1 -> AuthServer: Token Exchange Request\n(grant_type=urn:ietf:params:oauth:grant-type:token-exchange,\nrequested_token_type=urn:ietf:params:oauth:token-type:jwt,\nsubject_token=<Access Token>,\nsubject_token_type=urn:ietf:params:oauth:token-type:access_token,\nheaders=[{"Message-ID":<Message-ID>},\n{"From":"Alice <alice@example.com>"},\n{"Recipients-Digest":<Recipients-Digest>},\n{"Body-Digest":"nZURGvgk4xoy6-aI6dna5ddskq5ud_GyI7u96hkxYh4"}],\ndpop=<DPoP Proof>\nclient_id=<MTA1 client ID>,\nscope=internet_mail)
activate AuthServer

AuthServer -> AuthServer: Validate Access Token
AuthServer -> AuthServer: Validate DPoP Proof
AuthServer -> AuthServer: Verify that the "From" header (the email address of the sender)\n matches the email claim present in the Access Token
AuthServer -> AuthServer: Create JWT Assertion (iss=AS1,\n scope=internet_mail,\n azp=client_id,\n cnf={jkt=SHA256(JWK Thumbprint)},\n eh={eht=SHA256(Message-ID||From||Recipients-Digest||\n  Body-Digest),\n  ehl=["Message-ID","From","Recipients-Digest","Body-Digest"]})

AuthServer --> Entity1: Token Exchange Response\n(access_token=<JWT Assertion>,\ntoken_type=urn:ietf:params:oauth:token-type:jwt,\nissued_token_type=urn:ietf:params:oauth:token-type:jwt)
deactivate AuthServer

Entity1 --> Client1: The sending of the Placeholder Message\n has been authorized
deactivate Client1

Entity1 -> Entity1: Replace the stored Placeholder Message with the new JWT Signed Placeholder Message\n <Stored Placeholder Message>=<JWT Assertion>||<Placeholder Message>,\n use this Stored Placeholder Message as an access control list to grant or deny\n access to External Resources
Entity1 -> Entity1: Generate new DPoP Proof (using the same key pair, for MBX2/MTA2 service)

Entity1 -> Entity1: Use DNS to discover MBX2/MTA2 service using the domain part of recipient's\n email address
Entity1 -> Entity1: Oueue the Placeholder Message transfer request
activate Entity1 #LightGray

Entity1 -> Entity2: Transfer the Placeholder Message (message=<Placeholder Message>, Authorization: Bearer <JWT Assertion>, DPoP: <DPoP Proof>)
activate Entity2

Entity2 -> Entity2: Validate DPoP Proof (using "jkt" from JWT Assertion "cnf")
Entity2 -> Entity2: Validate JWT Assertion, use DNS discovery to verify that AS1\n is the corresponding JWT Assertion issuer
Entity2 -> Entity2: Verify the "eht" signature
Entity2 -> Entity2: Use DNS discovery to verify that MTA1 is the correct client\n making the request
Entity2 -> Entity2: Create JWT Signed Placeholder Message\n <Placeholder Message>=<JWT Assertion>||<message>
Entity2 -> Entity2: Create Received=<from, by, Date> header\n and add it to the Placeholder Message
Entity2 -> Entity2: Create Message-ID=<Message-ID> header\n and add it to the Placeholder Message
Entity2 -> Entity2: Store the Placeholder Message in the store,\n using the Message-ID header value as an identifier,\n <Stored Placeholder Message>=<Placeholder Message>
Entity2 --> Entity1: Placeholder Message Delivery Confirmation
'deactivate Entity2

Entity1 --> Entity1: Dequeue the Placeholder Message transfer request
deactivate Entity1
Entity1 -> Entity1: The Placeholder Message transfer request has been dequeued
deactivate Entity1

Entity2 -> Entity2: Queue the Placeholder Message proceeding using\n the Message-ID header value as an identifier
activate Entity2 #LightGray

@enduml
```

The sequence diagram of this flow is self-explanatory.

<div style="break-after:page"></div>

#### To authorize relaying/forwarding of the Placeholder Message:

* We will use the OAuth 2.0 Token Exchange grant type (RFC 8693).
* We will exchange the JWT Assertion for another JWT Assertion.
* We will use a JWT Assertion with the "eht" claim (a SHA-256 thumbprint of the set of Placeholder Message Envelope Headers that serves as a digital signature). This assertion embodies the DKIM signature.
* A forwarding policy can be set on the authorization server e.g., user may/may not be allowed to autoforward Placeholder Messages.
* The forwarding MBX2/MTA2 entity should not fetch the External Resources—only the final destination MBX3/MTA3 agent should.
* We will use the DPoP mechanism to bind the JWT Assertion to the client.
* The token exchange and relaying/forwarding process can begin after the successful Placeholder Message delivery.

Suppose we have another trust domain, example.org, and user Bob has two email addresses: bob@example.net and bobby@example.org. The sequence diagram below illustrates the relaying/forwarding process, in which the MBX2/MTA2 entity in the example.net trust domain forwards a Placeholder Message to the destination MBX3/MTA3 service operating within the example.org trust domain.

##### Sequence Diagram

```plantuml
@startuml
title OAuth 2.0 Token Exchange for relaying/forwarding Placeholder Message using JWT Assertion and DPoP
participant "MBX2/MTA2" as Entity2
participant "AS2 Authorization Server" as AuthServer
participant "MBX3/MTA3" as Entity3

activate Entity2
activate Entity2 #LightGray
Entity2 -> Entity2: Load the Placeholder Message using the Message-ID\n header value as an identifier from the store\n<Placeholder Message>=<Stored Placeholder Message>
Entity2 --> Entity2: The Placeholder Message proceeding dequeued
deactivate Entity2

Entity2 -> Entity2: Generate DPoP Proof (using new key pair)
Entity2 -> Entity2: Create JWT-Assertion-Digest=<SHA256(JWT Assertion)>\n header and add it to the Placeholder Message
Entity2 -> Entity2: Create Forwarded-From=<bob@example.net> header\n and add it to the Placeholder Message
Entity2 -> Entity2: Create Forwarded-To=<bobby@example.org> header\n and add it to the Placeholder Message
Entity2 -> Entity2: Create Recipients-Digest=<SHA256(Forwarded-To)> header\n and add it to the Placeholder Message
Entity2 -> Entity2: Create Date=<Date> header and add it to the Placeholder Message

Entity2 -> AuthServer: Token Exchange Request\n(grant_type=urn:ietf:params:oauth:grant-type:token-exchange,\nrequested_token_type=urn:ietf:params:oauth:token-type:jwt,\nsubject_token=<JWT Assertion>,\nsubject_token_type=urn:ietf:params:oauth:token-type:jwt,\nheaders=[{"Message-ID":<Message-ID>},\n{"Forwarded-From":"Bob <bob@example.net>"},\n{"Recipients-Digest":<Recipients-Digest>},\n{"Forwarded-To":"Bobby <bobby@example.org>},\n{"JWT-Assertion-Digest":<JWT-Assertion-Digest>}],\ndpop=<DPoP Proof>\nclient_id=<MTA2 client ID>\nscope=internet_mail)
activate AuthServer

AuthServer -> AuthServer: Validate JWT Assertion
AuthServer -> AuthServer: Validate DPoP Proof
AuthServer -> AuthServer: Evaluate the Forwarded-From header against\n the policy/rules set on the AS2
AuthServer -> AuthServer: Create JWT Assertion (iss=AS2,\n scope=internet_mail,\n azp=client_id,\n cnf={jkt=SHA256(JWK Thumbprint)},\n eh={eht=SHA256(Message-ID||Forwarded-From||\n  Recipients-Digest||JWT-Assertion-Digest),\n  ehl=["Message-ID","Forwarded-From","Recipients-Digest",\n  "JWT-Assertion-Digest"]})

AuthServer --> Entity2: Token Exchange Response\n(access_token=<JWT Assertion>,\ntoken_type=urn:ietf:params:oauth:token-type:jwt,\nissued_token_type=urn:ietf:params:oauth:token-type:jwt)
deactivate AuthServer

Entity2 -> Entity2: Replace the stored Placeholder Message with the new JWT Signed Placeholder Message\n <Stored Placeholder Message>=<JWT Assertion>||<Placeholder Message>,\n no need to use this Stored Placeholder Message as an access control list,\n the External Resources will be fetched from the origin MBX1/MTA1 service

Entity2 -> Entity2: Generate new DPoP Proof (using the same key pair, for MBX3/MTA3 service)

Entity2 -> Entity2: Use DNS to discover MBX3/MTA3 service using the domain part of\n recipient's email address

Entity2 -> Entity2: Oueue the Placeholder Message transfer request
activate Entity2 #LightGray

Entity2 -> Entity3: Transfer the Placeholder Message (message=<Placeholder Message>, Authorization: Bearer <JWT Assertion>, DPoP: <DPoP Proof>)
activate Entity3

Entity3 -> Entity3: Validate DPoP Proof (using "jkt" from JWT Assertion "cnf")
Entity3 -> Entity3: Validate JWT Assertion, use DNS discovery to verify that AS2\n is the corresponding JWT Assertion issuer
Entity3 -> Entity3: Verify the "eht" signature
Entity3 -> Entity3: Use DNS discovery to verify that MTA2 is the correct client\n making the request
Entity3 -> Entity3: Create JWT Signed Placeholder Message\n <Placeholder Message>=<JWT Assertion>||<message>
Entity3 -> Entity3: Create Received=<from, by, Date> header\n and add it to the Placeholder Message
Entity3 -> Entity3: Create Message-ID=<Message-ID> header\n and add it to the Placeholder Message
Entity3 -> Entity3: Store the Placeholder Message in the store using\n the Message-ID header value as an identifier
Entity3 --> Entity2: Placeholder Message Delivery Confirmation
'deactivate Entity3
Entity2 --> Entity2: Dequeue the Placeholder Message transfer request
deactivate Entity2
Entity2 -> Entity2: The Placeholder Message transfer request has been dequeued
deactivate Entity2

'Entity2 --> Entity1: Placeholder Message Forwarded Confirmation
'deactivate Entity2

Entity3 -> Entity3: Queue the Placeholder Message proceeding using\n the Message-ID header value as an identifier
activate Entity3 #LightGray

@enduml
```

The sequence diagram of this flow is self-explanatory. It outlines the process of forwarding the Placeholder Message to the destination MBX3/MTA3 service while chronologically appending the JWT Assertion to the top of the Placeholder Message.

<div style="break-after:page"></div>

#### To authorize the fetching of External Resources:

* We will use the OAuth 2.0 Token Exchange grant type (RFC 8693).
* We will exchange the JWT Assertion for an Access Token.
* We will use the DPoP mechanism to bind the Access Token to the client.
* The MBX3/MTA3 agent should be registered at the AS3 authorization server as a confidential client.
* The token exchange and fetching process begins after the successful Placeholder Message delivery.
* To fetch an External Resource, use the HTTP POST method to bypass caching. The request body should include the `envelope` parameter, which contains the (chain of) JWT Signed Placeholder Message Envelope Headers.

##### Sequence Diagram

```plantuml
@startuml
title OAuth 2.0 Token Exchange for fetching External Resources using JWT Assertion and DPoP
participant "MBX1/MTA1" as Entity1
participant "AS3 Authorization Server" as AuthServer
participant "MBX3/MTA3" as Entity3

activate Entity3
activate Entity3 #LightGray
Entity3 -> Entity3: Load the Placeholder Message using the Message-ID\n header value as an identifier from the store\n<Placeholder Message>=<Stored Placeholder Message>
Entity3 --> Entity3: The Placeholder Message proceeding dequeued
deactivate Entity3

Entity3 -> Entity3: Generate DPoP Proof (using new key pair)

Entity3 -> AuthServer: Token Exchange Request\n(grant_type=urn:ietf:params:oauth:grant-type:token-exchange,\naudience=MBX1/MTA1,\nrequested_token_type=urn:ietf:params:oauth:token-type:access_token,\nsubject_token=<JWT Assertion>,\nsubject_token_type=urn:ietf:params:oauth:token-type:jwt,\ndpop=<DPoP Proof>\nclient_id=<MTA3 client ID>,\nscope=internet_mail,\nmessage-id=<Message-ID>)
activate AuthServer

AuthServer -> AuthServer: Validate JWT Assertion
AuthServer -> AuthServer: Validate DPoP Proof
AuthServer -> AuthServer: Create Access Token (iss=AS3,\n exp=1737913610,\n aud=MBX1/MTA1,\n scope=internet_mail,\n azp=client_id,\n message-id=<Message-ID> ,\n cnf={jkt=SHA256(JWK Thumbprint)})

AuthServer --> Entity3: Token Exchange Response\n(access_token=<Access Token>,\ntoken_type=Bearer,\nissued_token_type=urn:ietf:params:oauth:token-type:access_token,\nexpires_in:3600)
deactivate AuthServer

Entity3 -> Entity3: Generate new DPoP Proof (using the same key pair,\n for MBX1/MTA1 service)

Entity3 -> Entity3: Use DNS to discover MBX1/MTA1 service using the domain part of\n sender's email address
Entity3 -> Entity3: Queue the External Resource fetching request
activate Entity3 #LightGray

Entity3 -> Entity1: Fetch External Resource (envelope=<Placeholder Message Envelope>, Content-ID, Authorization: Bearer <Access Token>, DPoP: <DPoP Proof>)
activate Entity1

Entity1 -> Entity1: Validate DPoP Proof (using "jkt" from Access Token "cnf")
Entity1 -> Entity1: Validate Access Token, use DNS discovery to verify that AS3 is the corresponding Access Token issuer
Entity1 -> Entity1: Use DNS discovery to verify that MTA3 is the correct client\n making the request
Entity1 -> Entity1: Load the Placeholder Message from the store using the message-id claim value from Access Token\n as an identifier
Entity1 -> Entity1: Verify the "envelope" parameter value against the Placeholder Message Envelope Headers
Entity1 -> Entity1: Optional, valid only for forwarded Placeholder Messages:\n Use the "envelope" parameter value from the fetching request to discharge the External Resources\n in accordance with the "Forwarded-From" and "Forwarded-To" headers
Entity1 -> Entity1: Use the Placeholder Message as an access control list to allow/deny access to the\n External Resource
Entity1 -> Entity1: Find External Resource in the store using Content-ID (file name of the External Resource equals\n the Content-ID header value)
Entity1 --> Entity3: Return External Resource
deactivate Entity1
Entity3 -> Entity3: Store the External Resource in the store using\n the Content-ID header value as an identifier
Entity3 --> Entity3: Dequeue the External Resource fetching request
deactivate Entity3
Entity3 -> Entity3: The External Resource fetching request has been dequeued

deactivate Entity3

@enduml
```

The sequence diagram of this flow is self-explanatory. Take note of the use of a JWT Signed Placeholder Message as an access control list during the fetch request.

<div style="break-after:page"></div>

## Project Implementation

We plan to implement this project as a proof of concept using Golang for the backend services, with a filesystem as the data store and SQLite for storing references and metadata of External Resources. OpenAPI 3.1 documentation will be developed to ensure clear and standardized API references. For the frontend, we will create a Progressive Web App (PWA) using the React framework, utilizing styled-components for styling. Additionally, we will design a synchronization protocol with polling to sync the Placeholder Message and inline External Resources with the mailbox.

## Future Work

1. Retry mechanism.

## Sketchy Ideas

1. Design a communication scheme to route the Placeholder Message through MBX/MTA intermediary entities, which will function as both proxy and reverse proxy services. This setup will facilitate the exchange of External Resources exclusively between the origin and the final destination. It is important to note that this scheme raises privacy considerations, as the intermediaries are not required to fetch the External Resources.
2. Email address verification link (verify email address, password reset) should be encrypted.
3. The user's public and private keys are stored in the AS. During an OAuth 2.0 redirect, the browser renders the user's private key in a hidden, read-only text field to enable cryptographic processes on the client only during user interaction.

# Prompt

Do you understand this project? Ask me if something is not clear to you.