function env () {
  console.log(process.env.REACT_APP_ENV)
  if (process.env.REACT_APP_ENV === 'development') {
    // 本地开发环境
    return {
      server: '//127.0.0.1:8088',
      wsServer: 'ws://127.0.0.1:8088',
      prefix: '',
    }
  } else {
    // 生产环境
    let wsPrefix;
    if (window.location.protocol === 'https:') {
      wsPrefix = 'wss:'
    } else {
      wsPrefix = 'ws:'
    }
    return {
      server: '',
      wsServer: wsPrefix + window.location.host,
      prefix: window.location.protocol + '//' + window.location.host,
    }
  }
}
export default env();

export const server = env().server;
export const wsServer = env().wsServer;
export const prefix = env().prefix;