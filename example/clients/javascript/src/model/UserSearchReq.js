/*
 * 
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * OpenAPI spec version: 
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 *
 * Swagger Codegen version: 2.4.18
 *
 * Do not edit the class manually.
 *
 */

(function(root, factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as an anonymous module.
    define(['ApiClient'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'));
  } else {
    // Browser globals (root is window)
    if (!root.SwaggerJsClient) {
      root.SwaggerJsClient = {};
    }
    root.SwaggerJsClient.UserSearchReq = factory(root.SwaggerJsClient.ApiClient);
  }
}(this, function(ApiClient) {
  'use strict';

  /**
   * The UserSearchReq model module.
   * @module model/UserSearchReq
   * @version 1.0.0
   */

  /**
   * Constructs a new <code>UserSearchReq</code>.
   * @alias module:model/UserSearchReq
   * @class
   */
  var exports = function() {
  };

  /**
   * Constructs a <code>UserSearchReq</code> from a plain JavaScript object, optionally creating a new instance.
   * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
   * @param {Object} data The plain JavaScript object bearing properties of interest.
   * @param {module:model/UserSearchReq} obj Optional instance to populate.
   * @return {module:model/UserSearchReq} The populated <code>UserSearchReq</code> instance.
   */
  exports.constructFromObject = function(data, obj) {
    if (data) {
      obj = obj || new exports();
      if (data.hasOwnProperty('keyWord'))
        obj.keyWord = ApiClient.convertToType(data['keyWord'], 'String');
    }
    return obj;
  }

  /**
   * @member {String} keyWord
   */
  exports.prototype.keyWord = undefined;


  return exports;

}));
