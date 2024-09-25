package Interfaces;

import java.util.concurrent.ConcurrentHashMap;

public interface Iwork {
    void add(Object key, Object value, int time);
    void delete(Object key);
    void put(Object key, Object value);
    ConcurrentHashMap<Object, Object> viewAll();
    String viewValueByKey(Object key);
    int size();
}
