package main;

import Interfaces.IWork;
import java.util.concurrent.ConcurrentHashMap;

public class Work implements IWork<String, String> {
    private final ConcurrentHashMap<String, String> concurrentHashMap = new ConcurrentHashMap<>();
    private final int time = 0;

    @Override
    public void add(String key, String value, int time) {
        this.concurrentHashMap.put(key, value);
    }

    @Override
    public void delete(String key) {
        this.concurrentHashMap.remove(key);
    }

    @Override
    public void put(String key, String value) {
        this.concurrentHashMap.replace(key, value);
    }

    @Override
    public String get(String key) {
        return this.concurrentHashMap.get(key);
    }

    @Override
    public int size(){
        return concurrentHashMap.size();
    }
}
