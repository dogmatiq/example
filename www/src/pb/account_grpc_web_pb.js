/**
 * @fileoverview gRPC-Web generated client stub for proto
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!



const grpc = {};
grpc.web = require('grpc-web');

const proto = {};
proto.proto = require('./account_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.proto.AccountClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

  /**
   * @private @const {?Object} The credentials to be used to connect
   *    to the server
   */
  this.credentials_ = credentials;

  /**
   * @private @const {?Object} Options for the client
   */
  this.options_ = options;
};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.proto.AccountPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

  /**
   * @private @const {!proto.proto.AccountClient} The delegate callback based client
   */
  this.delegateClient_ = new proto.proto.AccountClient(
      hostname, credentials, options);

};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.OpenAccountRequest,
 *   !proto.proto.OpenAccountResponse>}
 */
const methodInfo_Account_OpenAccount = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.OpenAccountResponse,
  /** @param {!proto.proto.OpenAccountRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.OpenAccountResponse.deserializeBinary
);


/**
 * @param {!proto.proto.OpenAccountRequest} request The
 *     request proto
 * @param {!Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.proto.OpenAccountResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.proto.OpenAccountResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.proto.AccountClient.prototype.openAccount =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/proto.Account/OpenAccount',
      request,
      metadata,
      methodInfo_Account_OpenAccount,
      callback);
};


/**
 * @param {!proto.proto.OpenAccountRequest} request The
 *     request proto
 * @param {!Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.proto.OpenAccountResponse>}
 *     The XHR Node Readable Stream
 */
proto.proto.AccountPromiseClient.prototype.openAccount =
    function(request, metadata) {
  return new Promise((resolve, reject) => {
    this.delegateClient_.openAccount(
      request, metadata, (error, response) => {
        error ? reject(error) : resolve(response);
      });
  });
};


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.proto.TestStreamingRequest,
 *   !proto.proto.TestStreamingResponse>}
 */
const methodInfo_Account_TestStreaming = new grpc.web.AbstractClientBase.MethodInfo(
  proto.proto.TestStreamingResponse,
  /** @param {!proto.proto.TestStreamingRequest} request */
  function(request) {
    return request.serializeBinary();
  },
  proto.proto.TestStreamingResponse.deserializeBinary
);


/**
 * @param {!proto.proto.TestStreamingRequest} request The request proto
 * @param {!Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.proto.TestStreamingResponse>}
 *     The XHR Node Readable Stream
 */
proto.proto.AccountClient.prototype.testStreaming =
    function(request, metadata) {
  return this.client_.serverStreaming(this.hostname_ +
      '/proto.Account/TestStreaming',
      request,
      metadata,
      methodInfo_Account_TestStreaming);
};


/**
 * @param {!proto.proto.TestStreamingRequest} request The request proto
 * @param {!Object<string, string>} metadata User defined
 *     call metadata
 * @return {!grpc.web.ClientReadableStream<!proto.proto.TestStreamingResponse>}
 *     The XHR Node Readable Stream
 */
proto.proto.AccountPromiseClient.prototype.testStreaming =
    function(request, metadata) {
  return this.delegateClient_.client_.serverStreaming(this.delegateClient_.hostname_ +
      '/proto.Account/TestStreaming',
      request,
      metadata,
      methodInfo_Account_TestStreaming);
};


module.exports = proto.proto;

