package main;

import Interfaces.Iwork;
import java.util.concurrent.ConcurrentHashMap;

public class workWithDB implements Iwork {
    private int time = 0;
    ConcurrentHashMap<Object, Object> concurrentHashMap = new ConcurrentHashMap<>();
    /*
     * TO DO:
     * Добавить возможность пункта со временем, чтобы в дальнейшем удалять ключи (ttl)
     * */
    @Override
    public void add(Object key, Object value, int time) {
        this.concurrentHashMap.put(key, value);
    }

    @Override
    public void delete(Object key) {
        this.concurrentHashMap.remove(key);
    }

    /*
    * TO DO:
    * Изменять время для ключа ttl
    * */
    @Override
    public void put(Object key, Object value) {
        if (this.concurrentHashMap.get(key) == null){
            this.concurrentHashMap.put(key, value);
        } else {
            this.concurrentHashMap.replace(key, value);
        }
    }

    @Override
    public ConcurrentHashMap<Object, Object> viewAll(){
        return concurrentHashMap;
    }

    @Override
    public String viewValueByKey(Object key) {
        return concurrentHashMap.get(key).toString();
    }

    @Override
    public int size(){
        return concurrentHashMap.size();
    }
}
