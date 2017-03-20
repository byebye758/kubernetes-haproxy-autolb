#!/usr/bin/python3
#-*- coding: utf-8 -*-

from jinja2 import Environment,  FileSystemLoader
from watchdog.events import FileSystemEventHandler
from watchdog.observers import Observer
import configparser
import socket, os, threading
import etcd3
import logging

#需要安装的包 etcd3,jinja2,watchdog,  python 3以上的版本

def getIp():
    #获取本机IP
    csock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    csock.connect(('114.114.114.114',53))
    (addr, port) = csock.getsockname()
    csock.close()
    return addr

class GetConfig():
    #获取配置文件
    def __init__(self,cnffile):
        cnf = configparser.ConfigParser()
        with open(cnffile,'r+') as cfgfile:
            cnf.read_file(cfgfile)
            self.etcdhost = cnf.get('global','etcdhost')
            self.etcdport = cnf.get('global','etcdport')
            self.etcdprefix = cnf.get('global','etcdprefix')
            self.templatename = cnf.get('global','templatename')
            self.templatepath = cnf.get('global','templatepath')
            self.hacnf = cnf.get('global','haproxyconfig')
            self.ha_lock_file = cnf.get('global','ha_lock_file')
            self.haproxy_lock_file = cnf.get('global','haproxy_lock_file')


class GetData(object):
    def __init__(self,**etcd):
        self.sourcedata = []
        self.project = []
        self.connetcd = etcd3.client(**etcd)

    def getprefix(self,gpath):
        #获取数据
        return self.connetcd.get_prefix(gpath)

    #监控
    #def watchprefix(self,wpath):
     #   self.connetcd.watch_prefix_once(wpath)

    def conversion_data(self, data):
        #把byte转化成字符串
        return eval(data[0].decode())

    def init_dict(self,podname,podip):
        #转化成字典
        return dict(Podname=podname,Podip=podip)

    def deldata(self,data):
        #删除无用的数据
        del data['Nodeip']
        del data['Haproxytable']
        del data['Lease']
        del data['Podip']
        del data['Podname']

    def dataprocess(self,ip,sd):
        #处理从etcd获取的数据
        for i in sd:
            d = self.conversion_data(i)
            if d['Haproxyip'] == ip:
                podlist = self.init_dict(d['Podname'], d['Podip'])
                d['Podlist'] = podlist
                self.deldata(d)
                if len(self.sourcedata) == 0:
                    # print(d['Projectname'])
                    d['Podlist'] = [podlist]
                    self.sourcedata.append(d)
                    self.project.append(d['Projectname'])
                else:
                    if d['Projectname'] not in self.project:
                        d['Podlist'] = [podlist]
                        self.sourcedata.append(d)
                        self.project.append(d['Projectname'])
                    elif d['Projectname'] in self.project:
                        for pn in self.sourcedata:
                            if pn['Projectname'] == d['Projectname']:
                                pn['Podlist'].append(d['Podlist'])
        return self.sourcedata


class MoBan():
    #模板生成配置文件
    def __init__(self,sourcefile,sourcepath,sourcedata,des_file):
        env = Environment(loader=FileSystemLoader(sourcepath), auto_reload=True)
        template = env.get_template(sourcefile)
        with open(des_file,'w+') as f:
            f.write(template.render(sourcedata))
            f.close()

class HaRegister():
    #启动ha注册程序
    def __init__(self,haproxy_lock_file,ha_lock_file):
        self.haproxypf = haproxy_lock_file
        self.hapf = ha_lock_file
    def checkHa(self):
        #检查haproxy是否运行，
        os.system('ps -ef|grep /usr/sbin/haproxy|grep -v grep > %s' % self.haproxypf)
        if not(os.path.getsize(self.haproxypf)):
            logging.getLogger('haconf').debug('haproxy is not running')
            os.system('systemctl start haproxy')
            logging.getLogger('haconf').debug('start haproxy')
        logging.getLogger('haconf').debug('haproxy is running')

    def starHaRegister(self):
        #检查ha注册程序是否运行
        os.system('ps -ef|grep /opt/haproxy/ha|grep -v grep > %s' % self.hapf)
        if not (os.path.getsize(self.hapf)):
            logging.getLogger('haconf').debug('ha register is not running')
            os.system('systemctl start ha')
            logging.getLogger('haconf').debug('start ha register')
        logging.getLogger('haconf').debug('ha register is running')


def updateHaConf(etcdhost, etcdport, etcdprefix, haip, hacnf,datalist):
    #根据模板生成配置文件
    data = datalist.getprefix(etcdprefix)
    msource = datalist.dataprocess(haip, data)
    logging.getLogger('haconf').debug(msource)
    MoBan(config.templatename, config.templatepath, {'moban': msource}, hacnf)

class Haproxy_Config_Handler(FileSystemEventHandler):
    #监控ha模板变化
    def on_modified(self, event):
        if event.src_path == "/opt/haproxy/templates/haproxy.template":
            logging.getLogger('haconf').debug('haproxy config template is change')
            sdata = GetData(host=config.etcdhost, port=config.etcdport)
            updateHaConf(config.etcdhost,config.etcdport,config.etcdprefix,getIp(),config.hacnf,sdata)
            #更新haproxy配置文件
            logging.getLogger('haconf').debug('update haproxy config is successful')
            os.system('systemctl reload haproxy')
            #动态加载haproxy
            logging.getLogger('haconf').debug('haproxy reload')


def configCheckTask():
    #检查模板的变化，当模板发生变化时重新更新一下配置文件
    logging.getLogger('haconf').info('start haproxy config template wtach')
    while True:
        event_handler = Haproxy_Config_Handler()
        observer = Observer()
        observer.schedule(event_handler,path=config.templatepath,recursive=False)
        try:
            observer.start()
        except Exception:
            logging.getLogger('haconf').debug('watch haproxy config template is failure')
            observer.stop()
            continue
        logging.getLogger('haconf').debug('update haproxy config is successful')
        observer.join()

def etcdWatchTask():
    #etcd监控
    logging.getLogger('haconf').info('start etcd watch')
    while True:
        logging.getLogger('haconf').debug('to get the data')
        sdata = GetData(host=config.etcdhost, port=config.etcdport)
        #获取etcd数据
        try:
            updateHaConf(config.etcdhost,config.etcdport,config.etcdprefix,getIp(),config.hacnf,sdata)
            #根据模板生成配置文件
        except Exception:
            #logging.getLogger('haconf').debug('update haproxy config is failure')
            sdata.connetcd.watch_prefix_once(config.etcdprefix)
            #监控etcd的变化
            logging.getLogger('haconf').debug('to get the data')
            continue
        logging.getLogger('haconf').debug('update haproxy config is successful')
        hareload = os.system('systemctl reload haproxy')
        #reload  haproxy
        if hareload:
           logging.getLogger('haconf').debug('haproxy reload is failure')
           continue
        logging.getLogger('haconf').debug('haproxy reload is successful')
        logging.getLogger('haconf').debug('start watch '+config.etcdprefix)
        sdata.connetcd.watch_prefix_once(config.etcdprefix)
        # 监控etcd的变化



if __name__ == '__main__':
    logging.basicConfig(format='%(asctime)s %(levelname)s %(name)s %(message)s',filename=os.path.join(os.getcwd(),'logs','haconf.log'),level=logging.DEBUG)
    logging.getLogger('haconf').info('load config file')
    #日志相关的配置,logging开头的都是打印日志
    config = GetConfig(os.path.join(os.getcwd(), 'hatemplate.cnf'))
    #加载配置文件

    HAR = HaRegister(config.haproxy_lock_file, config.ha_lock_file)
    #调用haproxy是否运行函数
    logging.getLogger('haconf').info('check haproxy process')
    HAR.checkHa()
    #调用ha注册程序是否运行函数
    logging.getLogger('haconf').info('check ha register process')
    HAR.starHaRegister()
    #启动多线程
    task1 = threading.Thread(target=configCheckTask)
    #调用模板文件监控函数
    task2 = threading.Thread(target=etcdWatchTask)
    #调用etcd监控函数

    task1.start()
    task2.start()

    task1.join()
    task2.join()