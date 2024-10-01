# Federizer

We introduce an enhancement to the current Internet mail architecture by incorporating authorization mechanisms for mail resources. The key components include:

* Mail Resource Store (MRS): A storage system where mail body resources are stored separately from the placeholder message.
* Placeholder Messages: Lightweight messages that contain references (URLs and cryptographic hashes) to the actual mail resources stored on the MRS.
* Resource Servers (RS) and Resource Retrieval Agents (RRA): Servers and agents responsible for storing and retrieving mail resources, with authentication mechanisms in place.
* Internet Mail Federation Protocol (Federizer): An open protocol designed to facilitate secure retrieval of mail resources between different security domains.

## Authorization-Enhanced Internet Mail Architecture

![Authorization-Enhanced Internet Mail Architecture](docs/src/main/images/authorization-enhanced_rfc5598.svg)

<p class="figure">
    Fig.&nbsp;1.&emsp;Authorization-enhanced Internet mail architecture using the Internet mail federation protocol.
</p>

Additional technical information: [Authorization-Enhanced Internet Mail Architecture](docs/Authorization-Enhanced_Internet_Mail_Architecture.md)

## Internet Mail Federation Protocol (Federizer)

Additional technical information: [Internet Mail Federation Protocol](docs/Internet_Mail_Federation_Protocol.md)

